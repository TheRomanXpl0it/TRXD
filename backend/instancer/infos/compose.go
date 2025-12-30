package infos

import (
	"encoding/json"
	"strconv"
)

type ComposeInfo struct {
	InstanceInfo
	ComposeBody string
	Env         map[string]string
}

func SetupComposeInfo(info *InstanceInfo, composeBody string) (*ComposeInfo, error) {
	composeInfo := ComposeInfo{
		InstanceInfo: *info,
		ComposeBody:  composeBody,
	}

	composeInfo.Env = make(map[string]string, 0)
	if info.Envs != "" {
		err := json.Unmarshal([]byte(info.Envs), &composeInfo.Env)
		if err != nil {
			return nil, err
		}
	}

	composeInfo.Env["MAX_MEMORY"] = strconv.Itoa(int(info.MaxMemory))
	composeInfo.Env["MAX_CPUS"] = info.MaxCpu
	composeInfo.Env["CONTAINER_NAME"] = info.Name
	composeInfo.Env["INSTANCE_DOMAIN"] = info.Domain
	if info.ExternalPort != nil {
		composeInfo.Env["INSTANCE_PORT"] = strconv.Itoa(int(*info.ExternalPort))
	}

	return &composeInfo, nil
}
