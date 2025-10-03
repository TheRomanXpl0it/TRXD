package infos

import (
	"encoding/json"
	"strconv"
	"strings"
)

type ContainerInfo struct {
	InstanceInfo
	Image           string
	Name            string
	Domain          string
	Env             []string
	ExternalPortStr string
	MaxCPUs         int64
}

func SetupContainerInfo(image string, info *InstanceInfo) (*ContainerInfo, error) {
	containerInfo := ContainerInfo{
		InstanceInfo: *info,
		Image:        image,
	}

	if info.ExternalPort == nil {
		splittedHost := strings.SplitN(info.Host, ".", 2)
		containerInfo.Name = splittedHost[0]
		containerInfo.Domain = splittedHost[1]
	} else {
		containerInfo.ExternalPortStr = strconv.Itoa(int(*info.ExternalPort))
		containerInfo.Name = containerInfo.ExternalPortStr
		containerInfo.Domain = info.Host
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

	containerInfo.Env = append(containerInfo.Env, "INSTANCE_HOST="+info.Host)
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
