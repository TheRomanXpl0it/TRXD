package instancer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"trxd/db/sqlc"
	"trxd/instancer/composes"
	"trxd/instancer/containers"
	"trxd/instancer/infos"

	"trxd/utils/log"
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

func spawnInstance(ctx context.Context, info *infos.InstanceInfo,
	conf *sqlc.GetDockerConfigsByIDRow, deployType sqlc.DeployType) (string, error) {

	var dockerID string
	var err error

	if info.UseDomain {
		info.NetID = "trxd-shared-internal"
	} else if deployType == sqlc.DeployTypeContainer {
		info.NetID = "trxd-shared-external"
	}

	if deployType == sqlc.DeployTypeContainer && conf.Image != "" {
		dockerID, err = containers.CreateContainer(ctx, info, conf.Image)
	} else if deployType == sqlc.DeployTypeCompose && conf.Compose != "" {
		dockerID, err = composes.CreateCompose(ctx, info, conf.Compose)
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

	log.Info("Creating instance:", "chall", challID, "team", tid)

	info, err := dbCreateInstance(ctx, tid, challID, expires_at, conf.HashDomain)
	if err != nil {
		return "", nil, err
	}
	if info == nil {
		cleanup = false
		return "", nil, errors.New("[race condition]")
	}

	var name string
	if !conf.HashDomain { // TODO: this can be omitted, use chall_cid_tid for all
		name = fmt.Sprintf("chall_%d_%d", challID, tid)
	} else {
		name = "chall_" + strings.Split(info.Host, ".")[0]
	}

	var externalPort *int32
	if info.Port.Valid {
		externalPort = &info.Port.Int32
	}

	hostRule := fmt.Sprintf("Host(`%s`)", info.Host)
	traefikRule := fmt.Sprintf("traefik.http.routers.%s.rule", name)
	traefikPort := fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", name)
	traefikEntrypoints := fmt.Sprintf("traefik.http.routers.%s.entrypoints", name)
	traefikPriority := fmt.Sprintf("traefik.http.routers.%s.priority", name)

	var labels map[string]string
	if conf.HashDomain {
		traefikPortValue := "1337"
		if internalPort != nil {
			traefikPortValue = fmt.Sprint(*internalPort)
		}

		labels = map[string]string{
			"traefik.enable":   "true",
			traefikRule:        hostRule,
			traefikPort:        traefikPortValue,
			traefikEntrypoints: "web",
			traefikPriority:    "10",
		}
	}

	instanceInfo := &infos.InstanceInfo{
		Name:         name,
		Domain:       info.Host,
		UseDomain:    conf.HashDomain,
		InternalPort: internalPort,
		ExternalPort: externalPort,
		Envs:         conf.Envs,
		MaxMemory:    int32(conf.MaxMemory.(int64)),
		MaxCpu:       conf.MaxCpu.(string),
		Labels:       labels,
	}

	dockerID, err = spawnInstance(ctx, instanceInfo, conf, deployType)
	if err != nil {
		return "", nil, err
	}

	err = dbUpdateInstanceDockerID(ctx, tid, challID, dockerID)
	if err != nil {
		return "", nil, err
	}

	cleanup = false

	return instanceInfo.Domain, instanceInfo.ExternalPort, nil
}
