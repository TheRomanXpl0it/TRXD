package instancer

import (
	"context"
	"errors"
	"fmt"
	"time"
	"trxd/db/sqlc"
	"trxd/instancer/composes"
	"trxd/instancer/containers"
	"trxd/instancer/infos"
	"trxd/instancer/networks"

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
	} else {
		netID, err := networks.CreateNetwork(ctx, instanceInfo.NetName)
		if err != nil {
			return "", nil, err
		}
		instanceInfo.NetID = netID
	}

	var res string
	if deployType == sqlc.DeployTypeContainer && conf.Image != "" {
		res, err = containers.CreateContainer(ctx, conf.Image, instanceInfo)
	} else if deployType == sqlc.DeployTypeCompose && conf.Compose != "" {
		projectName := fmt.Sprintf("chall_%d_%d", tid, challID)
		res, err = composes.CreateCompose(ctx, projectName, conf.Compose, instanceInfo)
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
