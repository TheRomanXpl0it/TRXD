package instancer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"trxd/db/sqlc"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/tde-nico/log"
)

type InstanceInfo struct {
	Host         string
	InternalPort *int32
	ExternalPort *int32
	Envs         string
	MaxMemory    int32
	MaxCpu       string
	NetName      string
	NetID        string
}

func FetchNginxID(ctx context.Context) (string, error) {
	args := filters.NewArgs()
	// args.Add("name", "trxd-nginx-1") // TODO: make dynamic
	args.Add("name", "trxd-test-pipeline-nginx-1") // TODO: make dynamic
	summary, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: args,
	})
	if err != nil {
		return "", err
	}
	if len(summary) != 1 {
		return "", fmt.Errorf("expected 1 network, got %d", len(summary))
	}

	return summary[0].ID, nil
}

func CreateInstance(ctx context.Context, tid, challID int32, internalPort *int32,
	expires_at time.Time, deployType sqlc.DeployType, conf *sqlc.GetDockerConfigsByIDRow) (string, *int32, error) {
	info, err := dbCreateInstance(ctx, tid, challID, expires_at, conf.HashDomain)
	if err != nil {
		return "", nil, err
	}
	if info == nil {
		return "", nil, errors.New("[race condition]")
	}

	log.Info("Creating instance:", "team", tid, "challenge", challID)

	instanceInfo := &InstanceInfo{
		Host:         info.Host,
		InternalPort: internalPort,
		Envs:         conf.Envs,
		MaxMemory:    int32(conf.MaxMemory.(int64)),
		MaxCpu:       conf.MaxCpu.(string),
		NetName:      fmt.Sprintf("net_%d_%d", tid, challID),
	}

	if !conf.HashDomain && info.Port.Valid {
		instanceInfo.ExternalPort = &info.Port.Int32
	} else if cli != nil {
		// TODO: already exists check
		net, err := cli.NetworkCreate(ctx, instanceInfo.NetName, network.CreateOptions{
			Internal: true,
		})
		if err != nil {
			return "", nil, err
		}

		instanceInfo.NetID = net.ID

		nginxID, err := FetchNginxID(ctx)
		if err != nil {
			return "", nil, err
		}

		err = cli.NetworkConnect(ctx, instanceInfo.NetID, nginxID, nil)
		if err != nil {
			return "", nil, err
		}
	}

	var res string
	if deployType == sqlc.DeployTypeContainer && conf.Image != "" {
		res, err = CreateContainer(ctx, conf.Image, instanceInfo)
	} else if deployType == sqlc.DeployTypeCompose && conf.Compose != "" {
		projectName := fmt.Sprintf("chall_%d_%d", tid, challID)
		res, err = CreateCompose(ctx, projectName, conf.Compose, instanceInfo)
	} else {
		return "", nil, errors.New("[no image or compose]")
	}
	if err != nil {
		return "", nil, err
	}

	err = UpdateInstanceDockerID(ctx, tid, challID, res)
	if err != nil {
		return "", nil, err
	}

	return instanceInfo.Host, instanceInfo.ExternalPort, nil
}

func CreateContainer(ctx context.Context, image string, info *InstanceInfo) (string, error) {
	if info.ExternalPort != nil && info.InternalPort == nil {
		return "", errors.New("[missing internal port]")
	}

	if cli == nil {
		return "", nil
	}

	var name, domain, portStr string
	if info.ExternalPort == nil {
		splittedHost := strings.SplitN(info.Host, ".", 2)
		name = splittedHost[0]
		domain = splittedHost[1]
	} else {
		portStr = strconv.Itoa(int(*info.ExternalPort))
		name = portStr
		domain = info.Host
	}

	env := make([]string, 0)
	if info.Envs != "" {
		var jsonEnvs map[string]string
		err := json.Unmarshal([]byte(info.Envs), &jsonEnvs)
		if err != nil {
			return "", err
		}
		for k, v := range jsonEnvs {
			env = append(env, k+"="+v)
		}
	}

	env = append(env, "INSTANCE_HOST="+info.Host)
	if portStr != "" {
		env = append(env, "INSTANCE_PORT="+portStr)
	}

	maxCPUs, err := strconv.ParseFloat(info.MaxCpu, 64)
	if err != nil {
		return "", err
	}

	containerConf := &container.Config{
		Hostname:     name,
		Domainname:   domain,
		Env:          env,
		Image:        image,
		ExposedPorts: nat.PortSet{},
	}
	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{},
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyAlways,
		},
		Resources: container.Resources{
			Memory:   int64(info.MaxMemory) * 1024 * 1024,
			NanoCPUs: int64(maxCPUs * 1e9),
		},
	}

	var networkingConfig *network.NetworkingConfig
	if portStr != "" {
		natPort := nat.Port(strconv.Itoa(int(*info.InternalPort)) + "/tcp")
		containerConf.ExposedPorts[natPort] = struct{}{}
		hostConf.PortBindings[natPort] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portStr}}
	} else {
		networkingConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				info.NetID: {},
			},
		}
	}

	if log.GetLevel() == log.DebugLevel {
		tmp1, err1 := json.MarshalIndent(containerConf, "", "  ")
		tmp2, err2 := json.MarshalIndent(hostConf, "", "  ")
		tmp3, err3 := json.MarshalIndent(networkingConfig, "", "  ")
		if err1 != nil || err2 != nil || err3 != nil {
			log.Debug("Created container:", "err", err1, "err", err2, "err", err3)
		} else {
			log.Debug("Created container:",
				"container", string(tmp1),
				"host", string(tmp2),
				"network", string(tmp3),
			)
		}
	}

	resp, err := cli.ContainerCreate(ctx, containerConf, hostConf, networkingConfig, nil, "chall_"+name)
	if err != nil {
		return "", err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func createComposeProject(ctx context.Context, name string, compose string, env map[string]string) (*types.Project, error) {
	configDetails := types.ConfigDetails{
		WorkingDir: "/" + name + "/",
		ConfigFiles: []types.ConfigFile{
			{Filename: "compose.yml", Content: []byte(compose)},
		},
		Environment: types.Mapping(env),
	}

	project, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(name, true)
	})
	if err != nil {
		return nil, err
	}

	for i, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False",
		}
		project.Services[i] = s
	}

	return project, nil
}

func CreateCompose(ctx context.Context, projectName string, composeBody string, info *InstanceInfo) (string, error) {
	if composeCli == nil {
		return "", nil
	}

	env := make(map[string]string, 0)
	if info.Envs != "" {
		err := json.Unmarshal([]byte(info.Envs), &env)
		if err != nil {
			return "", err
		}
	}

	env["MAX_MEMORY"] = strconv.Itoa(int(info.MaxMemory))
	env["MAX_CPUS"] = info.MaxCpu
	env["INSTANCE_HOST"] = info.Host
	if len(info.Host) > 0 {
		env["CONTAINER_NAME"] = "chall_" + strings.Split(info.Host, ".")[0]
	}
	if info.ExternalPort != nil {
		env["INSTANCE_PORT"] = strconv.Itoa(int(*info.ExternalPort))
	}

	project, err := createComposeProject(ctx, projectName, composeBody, env)
	if err != nil {
		return "", err
	}

	project.Networks["default"] = types.NetworkConfig{
		Name:     info.NetName,
		External: true,
	}

	if log.GetLevel() == log.DebugLevel {
		tmp, err := json.MarshalIndent(project, "", "  ")
		if err == nil {
			log.Debug("Created compose:", "project", string(tmp))
		} else {
			log.Debug("Created compose:", "err", err)
		}
	}

	err = composeCli.Up(ctx, project, api.UpOptions{})
	if err != nil {
		return "", err
	}

	return project.Name, nil
}
