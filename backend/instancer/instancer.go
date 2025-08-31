package instancer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/tde-nico/log"
)

var cli *client.Client
var composeCli api.Service

func InitInstancer() error {
	var err error

	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	composeCli, err = createComposeClient()
	if err != nil {
		return err
	}

	return nil
}

func createComposeClient() (api.Service, error) {
	dockerCli, err := command.NewDockerCli(command.WithAPIClient(cli))
	if err != nil {
		return nil, err
	}

	err = dockerCli.Initialize(&flags.ClientOptions{
		Context:  "default",
		LogLevel: "error",
	})
	if err != nil {
		return nil, err
	}

	srv := compose.NewComposeService(dockerCli)
	return srv, nil
}

func CreateInstance(ctx context.Context, tid, challID int32, expires_at time.Time,
	conf *sqlc.GetDockerConfigsByIDRow) (string, *int32, error) {
	info, err := dbCreateInstance(ctx, tid, challID, expires_at, conf.HashDomain)
	if err != nil {
		return "", nil, err
	}
	if info == nil { // race condition
		return "", nil, errors.New("[race condition]")
	}

	log.Info("Creating instance:", "team", tid, "challenge", challID)

	var port *int32
	if !conf.HashDomain && info.Port.Valid {
		port = &info.Port.Int32
	}
	if !conf.Envs.Valid {
		conf.Envs.String = ""
	}

	var res string
	if conf.Image.Valid {
		res, err = CreateContainer(ctx, conf.Image.String, info.Host, port,
			&conf.Envs.String, conf.MaxMemory.Int32, conf.MaxCpu.String)
	} else if conf.Compose.Valid {
		res, err = CreateCompose(ctx, conf.Compose.String, info.Host, port,
			&conf.Envs.String, conf.MaxMemory.Int32, conf.MaxCpu.String)
	} else {
		return "", nil, errors.New("[no image or compose]") // TODO: tests
	}
	if err != nil {
		return "", nil, err
	}

	err = UpdateInstanceDockerID(ctx, tid, challID, res)
	if err != nil {
		return "", nil, err
	}

	return info.Host, port, nil
}

func CreateContainer(ctx context.Context, image string, host string,
	port *int32, envs *string, maxMemory int32, maxCpu string) (string, error) {
	if cli == nil {
		return "", nil
	}

	var hash, domain string
	if port == nil {
		splittedHost := strings.SplitN(host, ".", 2)
		hash = splittedHost[0]
		domain = splittedHost[1]
		// TODO: add shared network
	} else {
		hash = strconv.Itoa(int(*port)) // TODO: use another variable
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
	if port != nil {
		env = append(env, "INSTANCE_PORT="+hash)
	}

	maxCPUs, err := strconv.ParseFloat(maxCpu, 64)
	if err != nil {
		return "", err
	}

	containerConf := &container.Config{
		Hostname:     hash,
		Domainname:   domain,
		Env:          env,
		Image:        image,
		ExposedPorts: nat.PortSet{},
	}
	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{},
		Resources: container.Resources{
			Memory:   int64(maxMemory) * 1024 * 1024,
			NanoCPUs: int64(maxCPUs * 1e9),
		},
	}

	if port != nil {
		natPort := nat.Port(hash + "/tcp")
		containerConf.ExposedPorts[natPort] = struct{}{} // TODO: change with the internal port
		hostConf.PortBindings[natPort] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hash}}
	}

	resp, err := cli.ContainerCreate(ctx, containerConf, hostConf, nil, nil, "chall_"+hash)
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

func CreateCompose(ctx context.Context, composeBody string, host string,
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

	env["INSTANCE_HOST"] = host
	env["MAX_MEMORY"] = strconv.Itoa(int(maxMemory))
	env["MAX_CPUS"] = maxCpu

	if port != nil {
		env["INSTANCE_PORT"] = strconv.Itoa(int(*port))
	}

	project, err := createComposeProject(ctx, "test-name", composeBody, env) // TODO: make real name
	if err != nil {
		return "", err
	}

	err = composeCli.Up(ctx, project, api.UpOptions{})
	if err != nil {
		return "", err
	}

	return project.Name, nil
}

func DeleteInstance(ctx context.Context, tid int32, challID int32, dockerID sql.NullString) error {
	log.Info("Deleting instance:", "team", tid, "challenge", challID)

	if dockerID.Valid {
		var err error
		log.Warn("Killing instance:", "docker_id", dockerID.String, "len", len(dockerID.String))
		if len(dockerID.String) == 64 { // TODO: maybe change this method of selection
			err = KillContainer(ctx, dockerID.String)
		} else {
			err = KillCompose(ctx, dockerID.String)
		}
		if err != nil {
			return err
		}
	}

	err := dbDeleteInstance(ctx, tid, challID)
	if err != nil {
		return err
	}

	return nil
}

func KillContainer(ctx context.Context, id string) error {
	if cli == nil {
		return nil
	}

	err := cli.ContainerKill(ctx, id, "SIGKILL")
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, id, container.RemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}

func KillCompose(ctx context.Context, name string) error {
	if composeCli == nil {
		return nil
	}

	err := composeCli.Down(ctx, name, api.DownOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetInterval() (time.Duration, error) {
	ctx := context.Background()
	conf, err := db.GetConfig(ctx, "reclaim-instance-interval")
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(conf.Value)
	if err != nil {
		return 0, err
	}

	sleep := time.Duration(value) * time.Second
	return sleep, nil
}

func ReclaimLoop() {
	err := InitInstancer()
	if err != nil {
		log.Fatal("Failed to initialize instancer:", "err", err)
	}
	defer cli.Close()

	sleep, err := GetInterval()
	if err != nil {
		log.Fatal("Failed to get reclaim interval:", "err", err)
	}

	for {
		time.Sleep(sleep)
		ctx := context.Background()

		next, err := db.Sql.GetNextInstanceToDelete(ctx)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Failed to get next instance to delete:", "err", err)
			} else {
				sleep, err = GetInterval()
				if err != nil {
					log.Fatal("Failed to get reclaim interval:", "err", err)
				}
			}
			continue
		}

		if time.Now().Before(next.ExpiresAt) {
			sleep = time.Until(next.ExpiresAt)
			continue
		} else {
			sleep = 0
		}

		err = DeleteInstance(ctx, next.TeamID, next.ChallID, next.DockerID)
		if err != nil {
			log.Error("Failed to delete instance:", "err", err)
		}
	}
}
