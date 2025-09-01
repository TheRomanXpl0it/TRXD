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
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/tde-nico/log"
)

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

	var externalPort *int32
	if !conf.HashDomain && info.Port.Valid {
		externalPort = &info.Port.Int32
	}
	if !conf.Envs.Valid {
		conf.Envs.String = ""
	}

	var res string
	if deployType == sqlc.DeployTypeContainer && conf.Image.Valid {
		res, err = CreateContainer(ctx, conf.Image.String, info.Host, internalPort,
			externalPort, &conf.Envs.String, conf.MaxMemory.Int32, conf.MaxCpu.String)
	} else if deployType == sqlc.DeployTypeCompose && conf.Compose.Valid {
		projectName := fmt.Sprintf("chall_%d_%d", tid, challID)
		res, err = CreateCompose(ctx, projectName, conf.Compose.String, info.Host, externalPort,
			&conf.Envs.String, conf.MaxMemory.Int32, conf.MaxCpu.String)
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

	return info.Host, externalPort, nil
}

func CreateContainer(ctx context.Context, image string, host string, internalPort *int32,
	externalPort *int32, envs *string, maxMemory int32, maxCpu string) (string, error) {
	if externalPort != nil && internalPort == nil {
		return "", errors.New("[missing internal port]")
	}

	if cli == nil {
		return "", nil
	}

	var name, domain, portStr string
	if externalPort == nil {
		splittedHost := strings.SplitN(host, ".", 2)
		name = splittedHost[0]
		domain = splittedHost[1]
	} else {
		portStr = strconv.Itoa(int(*externalPort))
		name = portStr
		domain = host
	}

	env := make([]string, 0)
	if envs != nil && *envs != "" {
		var jsonEnvs map[string]string
		err := json.Unmarshal([]byte(*envs), &jsonEnvs)
		if err != nil {
			return "", err
		}
		for k, v := range jsonEnvs {
			env = append(env, k+"="+v)
		}
	}

	env = append(env, "INSTANCE_HOST="+host)
	if portStr != "" {
		env = append(env, "INSTANCE_PORT="+portStr)
	}

	maxCPUs, err := strconv.ParseFloat(maxCpu, 64)
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
			Memory:   int64(maxMemory) * 1024 * 1024,
			NanoCPUs: int64(maxCPUs * 1e9),
		},
	}

	var networkingConfig *network.NetworkingConfig
	if portStr != "" {
		natPort := nat.Port(strconv.Itoa(int(*internalPort)) + "/tcp")
		containerConf.ExposedPorts[natPort] = struct{}{}
		hostConf.PortBindings[natPort] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: portStr}}
	} else {
		networkingConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				SHARED_NETWORK: {},
			},
		}
	}

	if log.GetLevel() == log.DebugLevel {
		tmp1, err1 := json.MarshalIndent(containerConf, "", "  ")
		tmp2, err2 := json.MarshalIndent(hostConf, "", "  ")
		tmp3, err3 := json.MarshalIndent(networkingConfig, "", "  ")
		if err1 == nil && err2 == nil && err3 == nil {
			log.Debug("Created container:",
				"container", string(tmp1),
				"host", string(tmp2),
				"network", string(tmp3),
			)
		} else {
			log.Debug("Created container:", "err", err1, "err", err2, "err", err3)
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

func CreateCompose(ctx context.Context, projectName string, composeBody string, host string,
	port *int32, envs *string, maxMemory int32, maxCpu string) (string, error) {
	if composeCli == nil {
		return "", nil
	}

	env := make(map[string]string, 0)
	if envs != nil && *envs != "" {
		err := json.Unmarshal([]byte(*envs), &env)
		if err != nil {
			return "", err
		}
	}

	env["MAX_MEMORY"] = strconv.Itoa(int(maxMemory))
	env["MAX_CPUS"] = maxCpu
	env["INSTANCE_HOST"] = host
	if len(host) > 0 {
		env["CONTAINER_NAME"] = "chall_" + strings.Split(host, ".")[0]
	}
	if port != nil {
		env["INSTANCE_PORT"] = strconv.Itoa(int(*port))
	}

	project, err := createComposeProject(ctx, projectName, composeBody, env)
	if err != nil {
		return "", err
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
