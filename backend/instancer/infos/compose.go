package infos

import (
	"encoding/json"
	"strconv"
	"strings"
)

type ComposeInfo struct {
	InstanceInfo
	ProjectName string
	ComposeBody string
	Env         map[string]string
}

func SetupComposeInfo(projectName string, composeBody string, info *InstanceInfo) (*ComposeInfo, error) {
	composeInfo := ComposeInfo{
		InstanceInfo: *info,
		ProjectName:  projectName,
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
	composeInfo.Env["INSTANCE_HOST"] = info.Host
	if len(info.Host) > 0 {
		composeInfo.Env["CONTAINER_NAME"] = "chall_" + strings.Split(info.Host, ".")[0]
	} else {
		composeInfo.Env["CONTAINER_NAME"] = projectName
	}
	if info.ExternalPort != nil {
		composeInfo.Env["INSTANCE_PORT"] = strconv.Itoa(int(*info.ExternalPort))
	}

	return &composeInfo, nil
}
