package instancer

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"trxd/api/routes/instances_delete"
	"trxd/db"

	"github.com/moby/moby/client"
	"github.com/tde-nico/log"
)

// TODO

type Instancer struct {
	cli *client.Client
}

// TODO: defer cli.Close()

func CreateInstancer() (*Instancer, error) {
	i := &Instancer{}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	i.cli = cli

	return i, nil
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

		err = instances_delete.DeleteInstance(ctx, next.TeamID, next.ChallID)
		if err != nil {
			log.Error("Failed to delete instance:", "err", err)
		}

		// TODO: delete the instance
		log.Info("Reclaiming instance: team", next.TeamID, "challenge", next.ChallID)
	}
}
