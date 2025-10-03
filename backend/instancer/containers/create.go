package containers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"trxd/instancer/infos"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/tde-nico/log"
)

func CreateContainer(ctx context.Context, image string, info *infos.InstanceInfo) (string, error) {
	if info.ExternalPort != nil && info.InternalPort == nil {
		return "", errors.New("[missing internal port]")
	}

	if Cli == nil {
		return "", nil
	}

	containerInfo, err := infos.SetupContainerInfo(image, info)
	if err != nil {
		return "", err
	}

	containerConf, hostConf, networkingConfig, err := setupContainerConf(containerInfo)
	if err != nil {
		return "", err
	}

	if log.GetLevel() == log.DebugLevel {
		debugContainer(containerConf, hostConf, networkingConfig)
	}

	var containerID string
	resp, err := Cli.ContainerCreate(ctx, containerConf, hostConf, networkingConfig, nil, "chall_"+containerInfo.Name)
	if err == nil {
		containerID = resp.ID
	} else {
		if !strings.Contains(err.Error(), "is already in use") {
			return "", err
		}

		containerID, err = FetchContainerByName(ctx, "chall_"+containerInfo.Name)
		if err != nil {
			return "", err
		}

		if info.NetID != "" {
			err = Cli.NetworkConnect(ctx, info.NetID, containerID, nil)
			if err != nil {
				if !strings.Contains(err.Error(), "already exists in network") {
					return "", err
				}
			}
		}
	}

	err = Cli.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	return containerID, nil
}

func setupContainerConf(info *infos.ContainerInfo) (*container.Config, *container.HostConfig, *network.NetworkingConfig, error) {
	containerConf := &container.Config{
		Hostname:     info.Name,
		Domainname:   info.Domain,
		Env:          info.Env,
		Image:        info.Image,
		ExposedPorts: nat.PortSet{},
	}

	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{},
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyAlways,
		},
		Resources: container.Resources{
			Memory:   int64(info.MaxMemory) * 1024 * 1024,
			NanoCPUs: info.MaxCPUs,
		},
	}

	var networkingConfig *network.NetworkingConfig
	if info.ExternalPortStr != "" {
		natPort := nat.Port(strconv.Itoa(int(*info.InternalPort)) + "/tcp")
		containerConf.ExposedPorts[natPort] = struct{}{}
		hostConf.PortBindings[natPort] = []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: info.ExternalPortStr,
		}}
	} else {
		networkingConfig = &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				info.NetID: {},
			},
		}
	}

	return containerConf, hostConf, networkingConfig, nil
}

func debugContainer(containerConf *container.Config, hostConf *container.HostConfig, networkingConfig *network.NetworkingConfig) {
	tmp1, err1 := json.MarshalIndent(containerConf, "", "  ")
	tmp2, err2 := json.MarshalIndent(hostConf, "", "  ")
	tmp3, err3 := json.MarshalIndent(networkingConfig, "", "  ")
	if err1 != nil || err2 != nil || err3 != nil {
		log.Error("Created container:", "err", err1, "err", err2, "err", err3)
	} else {
		log.Debug("Created container:",
			"container", string(tmp1),
			"host", string(tmp2),
			"network", string(tmp3),
		)
	}
}
