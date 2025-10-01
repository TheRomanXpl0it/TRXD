package instancer

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"trxd/db"
	"trxd/utils/consts"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
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

	dockerCli, err := command.NewDockerCli(command.WithAPIClient(cli))
	if err != nil {
		return err
	}

	err = dockerCli.Initialize(&flags.ClientOptions{
		Context:  "default",
		LogLevel: "error",
	})
	if err != nil {
		return err
	}

	composeCli = compose.NewComposeService(dockerCli)

	return nil
}

func GetInterval() (time.Duration, error) {
	ctx := context.Background()
	conf, err := db.GetConfig(ctx, "reclaim-instance-interval")
	if err != nil {
		return 0, err
	}
	if conf == "" {
		interval := consts.DefaultConfigs["reclaim-instance-interval"].(int)
		return time.Duration(interval) * time.Second, nil
	}

	value, err := strconv.Atoi(conf)
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
