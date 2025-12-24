package instancer

import (
	"context"
	"database/sql"
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

func recoverBrokenInstance(ctx context.Context, tid int32, challID int32, dockerID string) {
	err := DeleteInstance(ctx, tid, challID, sql.NullString{String: dockerID, Valid: dockerID != ""})
	if err == nil {
		return
	}
	log.Error("Failed to cleanup instance after creation failure", "team", tid, "challenge", challID, "err", err)

	err = UpdateInstanceExpire(ctx, tid, challID, time.Now().Add(-1*time.Second))
	if err == nil {
		return
	}

	log.Error("Failed to expire instance after creation failure", "team", tid, "challenge", challID, "err", err)
}

func spawnInstance(ctx context.Context, instanceInfo *infos.InstanceInfo,
	conf *sqlc.GetDockerConfigsByIDRow, deployType sqlc.DeployType, port *int32, name string) (string, error) {

	var dockerID string
	var err error

	if !conf.HashDomain && port != nil {
		instanceInfo.ExternalPort = port
		if deployType == sqlc.DeployTypeContainer {
			instanceInfo.NetID = "trxd-shared"
		}
	} else {
		netID, err := networks.CreateNetwork(ctx, instanceInfo.NetName, false)
		if err != nil {
			return "", err
		}
		instanceInfo.NetID = netID
	}

	if deployType == sqlc.DeployTypeContainer && conf.Image != "" {
		dockerID, err = containers.CreateContainer(ctx, name, conf.Image, instanceInfo)
	} else if deployType == sqlc.DeployTypeCompose && conf.Compose != "" {
		dockerID, err = composes.CreateCompose(ctx, name, conf.Compose, instanceInfo)
	} else {
		return "", errors.New("[no image or compose]")
	}
	if err != nil {
		return dockerID, err
	}

	return dockerID, nil
}

func CreateInstance(ctx context.Context, tid int32, challID int32, internalPort *int32, expires_at time.Time,
	deployType sqlc.DeployType, conf *sqlc.GetDockerConfigsByIDRow) (string, *int32, error) {

	var dockerID string
	cleanup := true

	defer func() {
		r := recover()
		if r == nil && !cleanup {
			return
		}
		if r != nil {
			log.Critical("Recovered instancer create panic", "crit", r)
		}

		recoverBrokenInstance(ctx, tid, challID, dockerID)
	}()

	log.Info("Creating instance:", "team", tid, "challenge", challID)

	info, err := dbCreateInstance(ctx, tid, challID, expires_at, conf.HashDomain)
	if err != nil {
		return "", nil, err
	}
	if info == nil {
		cleanup = false
		return "", nil, errors.New("[race condition]")
	}

	instanceInfo := &infos.InstanceInfo{
		Host:         info.Host,
		InternalPort: internalPort,
		Envs:         conf.Envs,
		MaxMemory:    int32(conf.MaxMemory.(int64)),
		MaxCpu:       conf.MaxCpu.(string),
		NetName:      fmt.Sprintf("net_%d_%d", challID, tid),
	}
	var port *int32
	if info.Port.Valid {
		port = &info.Port.Int32
	}
	name := fmt.Sprintf("chall_%d_%d", challID, tid)

	dockerID, err = spawnInstance(ctx, instanceInfo, conf, deployType, port, name)
	if err != nil {
		return "", nil, err
	}

	err = dbUpdateInstanceDockerID(ctx, tid, challID, dockerID)
	if err != nil {
		return "", nil, err
	}

	cleanup = false

	return instanceInfo.Host, instanceInfo.ExternalPort, nil
}
