package instancer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/instancer/composes"
	"trxd/instancer/containers"
	"trxd/instancer/infos"
	"trxd/instancer/networks"

	"github.com/tde-nico/log"
)

func CreateInstance(ctx context.Context, tid, challID int32, internalPort *int32, expires_at time.Time,
	deployType sqlc.DeployType, conf *sqlc.GetDockerConfigsByIDRow) (string, *int32, error) {
	var dockerID string
	success := false
	inst_lock := fmt.Sprintf("inst_lock_%d_%d", tid, challID)

	status, err := db.StorageSetNX(ctx, inst_lock, "true")
	if err != nil {
		return "", nil, err
	}
	if !status {
		return "", nil, errors.New("[race condition]")
	}

	defer func() {
		err := db.StorageDelete(ctx, inst_lock)
		if err != nil {
			log.Error("Failed to release instance lock", "team", tid, "challenge", challID, "err", err)
		}

		if success {
			return
		}

		err = DeleteInstance(ctx, tid, challID, sql.NullString{String: dockerID, Valid: dockerID != ""})
		if err == nil {
			return
		}
		log.Error("Failed to cleanup instance after creation failure", "team", tid, "challenge", challID, "err", err)

		err = UpdateInstanceExpire(ctx, tid, challID, time.Now().Add(-1*time.Second))
		if err == nil {
			return
		}
		log.Error("Failed to expire instance after creation failure", "team", tid, "challenge", challID, "err", err)
	}()

	log.Info("Creating instance:", "team", tid, "challenge", challID)

	info, err := dbCreateInstance(ctx, tid, challID, expires_at, conf.HashDomain)
	if err != nil {
		return "", nil, err
	}
	if info == nil {
		return "", nil, errors.New("[race condition]")
	}

	instanceInfo := &infos.InstanceInfo{
		Host:         info.Host,
		InternalPort: internalPort,
		Envs:         conf.Envs,
		MaxMemory:    int32(conf.MaxMemory.(int64)),
		MaxCpu:       conf.MaxCpu.(string),
		NetName:      fmt.Sprintf("net_%d_%d", tid, challID),
	}

	if !conf.HashDomain && info.Port.Valid {
		instanceInfo.ExternalPort = &info.Port.Int32
		if deployType == sqlc.DeployTypeContainer {
			instanceInfo.NetID = "trxd-shared"
		}
	} else {
		netID, err := networks.CreateNetwork(ctx, instanceInfo.NetName, false)
		if err != nil {
			return "", nil, err
		}
		instanceInfo.NetID = netID
	}

	if deployType == sqlc.DeployTypeContainer && conf.Image != "" {
		dockerID, err = containers.CreateContainer(ctx, conf.Image, instanceInfo)
	} else if deployType == sqlc.DeployTypeCompose && conf.Compose != "" {
		projectName := fmt.Sprintf("chall_%d_%d", tid, challID)
		dockerID, err = composes.CreateCompose(ctx, projectName, conf.Compose, instanceInfo)
	} else {
		return "", nil, errors.New("[no image or compose]")
	}
	if err != nil {
		return "", nil, err
	}

	err = dbUpdateInstanceDockerID(ctx, tid, challID, dockerID)
	if err != nil {
		return "", nil, err
	}

	success = true

	return instanceInfo.Host, instanceInfo.ExternalPort, nil
}
