package infos

import (
	"encoding/json"
	"strconv"
)

type ContainerInfo struct {
	InstanceInfo
	Image           string
	ExternalPortStr string
	Env             []string
	MaxCPUs         int64
}

func SetupContainerInfo(info *InstanceInfo, image string) (*ContainerInfo, error) {
	containerInfo := ContainerInfo{
		InstanceInfo: *info,
		Image:        image,
	}

	if info.ExternalPort != nil {
		containerInfo.ExternalPortStr = strconv.Itoa(int(*info.ExternalPort))
	}

	containerInfo.Env = make([]string, 0)
	if info.Envs != "" {
		var jsonEnvs map[string]string
		err := json.Unmarshal([]byte(info.Envs), &jsonEnvs)
		if err != nil {
			return nil, err
		}
		for k, v := range jsonEnvs {
			containerInfo.Env = append(containerInfo.Env, k+"="+v)
		}
	}

	containerInfo.Env = append(containerInfo.Env, "INSTANCE_DOMAIN="+info.Domain)
	if containerInfo.ExternalPortStr != "" {
		containerInfo.Env = append(containerInfo.Env, "INSTANCE_PORT="+containerInfo.ExternalPortStr)
	}

	maxCPUs, err := strconv.ParseFloat(info.MaxCpu, 64)
	if err != nil {
		return nil, err
	}
	containerInfo.MaxCPUs = int64(maxCPUs * 1e9)

	return &containerInfo, nil
}
