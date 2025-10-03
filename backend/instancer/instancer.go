package instancer

import (
	"context"
	"database/sql"
	"strconv"
	"time"
	"trxd/db"
	"trxd/instancer/composes"
	"trxd/instancer/containers"
	"trxd/utils/consts"

	"github.com/tde-nico/log"
)

func InitInstancer() error {
	var err error

	err = containers.InitCli()
	if err != nil {
		return err
	}

	err = composes.InitComposeCli()
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
	defer containers.CloseCli()

	// TODO: make this a separate function with anti panic

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
