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

	"trxd/utils/consts"
	"trxd/utils/log"
)

type CreateInstanceParams struct {
	Tid          int32
	ChallID      int32
	ConnType     sqlc.ConnType
	InternalPort *int32
	DeployType   sqlc.DeployType
	DockerConfig *sqlc.GetDockerConfigsByIDRow
}

type CreateInstanceResult struct {
	Host       string
	Port       *int32
	Expiration time.Time
}

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

func makeLabels(info *infos.InstanceInfo, p *CreateInstanceParams) {
	if !p.DockerConfig.HashDomain {
		return
	}

	routersRule := "traefik.%s.routers.%s.rule"
	routersEntrypoints := "traefik.%s.routers.%s.entrypoints"
	routersTls := "traefik.%s.routers.%s.tls"
	routersPriotity := "traefik.%s.routers.%s.priority"
	loadbalancerPort := "traefik.%s.services.%s.loadbalancer.server.port"

	var protocol, rule, entrypoint string
	if p.ConnType == sqlc.ConnTypeTCP {
		protocol = "tcp"
		rule = "HostSNI(`%s`)"
		entrypoint = "tcp"
	} else { // so http is (NONE, HTTP, HTTPS)
		protocol = "http"
		rule = "Host(`%s`)"
		entrypoint = "web"
	}

	traefikPort := "1337"
	if p.InternalPort != nil {
		traefikPort = fmt.Sprint(*p.InternalPort)
	}

	traefikRoutersRule := fmt.Sprintf(routersRule, protocol, info.Name)
	traefikRoutersEntrypoints := fmt.Sprintf(routersEntrypoints, protocol, info.Name)
	traefikRoutersPriority := fmt.Sprintf(routersPriotity, protocol, info.Name)
	traefikLoadbalancerPort := fmt.Sprintf(loadbalancerPort, protocol, info.Name)
	traefikRoutersTls := fmt.Sprintf(routersTls, protocol, info.Name)

	info.Labels = map[string]string{
		"traefik.enable":          "true",
		"traefik.docker.network":  consts.NetworkInternal,
		traefikRoutersRule:        fmt.Sprintf(rule, info.Domain),
		traefikRoutersEntrypoints: entrypoint,
		traefikRoutersPriority:    "10",
		traefikLoadbalancerPort:   traefikPort,
	}

	if protocol == "tcp" {
		info.Labels[traefikRoutersTls] = "true"
	}
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

func CreateInstance(ctx context.Context, p *CreateInstanceParams) (*CreateInstanceResult, error) {
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

		recoverBrokenInstance(ctx, p.Tid, p.ChallID, dockerID)
	}()

	log.Info("Creating instance:", "chall", p.ChallID, "team", p.Tid)

	lifetime := time.Second * time.Duration(p.DockerConfig.Lifetime.(int64))
	expires_at := time.Now().Add(lifetime)

	creationInfo, err := dbCreateInstance(ctx, p.Tid, p.ChallID, expires_at, p.DockerConfig.HashDomain)
	if err != nil {
		return nil, err
	}
	if creationInfo == nil {
		cleanup = false
		return nil, errors.New("[race condition]")
	}

	instanceInfo := &infos.InstanceInfo{
		Name:         fmt.Sprintf("chall_%d_%d", p.ChallID, p.Tid),
		Domain:       creationInfo.Host,
		UseDomain:    p.DockerConfig.HashDomain,
		InternalPort: p.InternalPort,
		Envs:         p.DockerConfig.Envs,
		MaxMemory:    int32(p.DockerConfig.MaxMemory.(int64)),
		MaxCpu:       p.DockerConfig.MaxCpu.(string),
	}

	if creationInfo.Port.Valid {
		instanceInfo.ExternalPort = &creationInfo.Port.Int32
	}

	makeLabels(instanceInfo, p)

	dockerID, err = spawnInstance(ctx, instanceInfo, p.DockerConfig, p.DeployType)
	if err != nil {
		return nil, err
	}

	err = dbUpdateInstanceDockerID(ctx, p.Tid, p.ChallID, dockerID)
	if err != nil {
		return nil, err
	}

	cleanup = false

	return &CreateInstanceResult{
		Host:       instanceInfo.Domain,
		Port:       instanceInfo.ExternalPort,
		Expiration: expires_at,
	}, nil
}
