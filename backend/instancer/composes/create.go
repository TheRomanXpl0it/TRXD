package composes

import (
	"context"
	"encoding/json"
	"strings"
	"trxd/instancer/infos"

	"trxd/utils/log"

	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v5/pkg/api"
)

func CreateCompose(ctx context.Context, info *infos.InstanceInfo, composeBody string) (string, error) {
	if ComposeCli == nil {
		return "", nil
	}

	composeInfo, err := infos.SetupComposeInfo(info, composeBody)
	if err != nil {
		return "", err
	}

	project, err := setupComposeProject(ctx, composeInfo)
	if err != nil {
		return "", err
	}

	if log.GetLevel() == log.DebugLevel {
		debugCompose(project)
	}

	err = ComposeCli.Up(ctx, project, api.UpOptions{})
	if err != nil {
		return "", err
	}

	return project.Name, nil
}

func setupComposeProject(ctx context.Context, info *infos.ComposeInfo) (*types.Project, error) {
	configDetails := types.ConfigDetails{
		WorkingDir: "/" + info.Name + "/",
		ConfigFiles: []types.ConfigFile{
			{Filename: "compose.yml", Content: []byte(info.ComposeBody)},
		},
		Environment: types.Mapping(info.Env),
	}

	project, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		options.SetProjectName(info.Name, true)
	})
	if err != nil {
		return nil, err
	}

	for i, s := range project.Services {
		s.CustomLabels = map[string]string{
			api.ProjectLabel:     project.Name,
			api.ServiceLabel:     s.Name,
			api.VersionLabel:     api.ComposeVersion,
			api.WorkingDirLabel:  "/",
			api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
			api.OneoffLabel:      "False",
		}

		if s.Name == "chall" {
			for k, v := range info.Labels {
				s.CustomLabels[k] = v
			}
		}

		project.Services[i] = s
	}

	if info.NetID != "" {
		project.Networks["default"] = types.NetworkConfig{
			Name:     info.NetName,
			External: true,
		}
	}

	return project, nil
}

func debugCompose(project *types.Project) {
	tmp, err := json.MarshalIndent(project, "", "  ")
	if err == nil {
		log.Debug("Created compose:", "project", string(tmp))
	} else {
		log.Debug("Created compose:", "err", err)
	}
}
