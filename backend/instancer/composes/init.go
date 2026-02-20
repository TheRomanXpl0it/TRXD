package composes

import (
	"trxd/instancer/containers"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

var ComposeCli api.Compose

func InitComposeCli() error {
	if containers.Cli == nil {
		return nil
	}

	dockerCli, err := command.NewDockerCli(command.WithAPIClient(containers.Cli))
	if err != nil {
		return err
	}

	err = dockerCli.Initialize(&flags.ClientOptions{
		Context:  "default",
		LogLevel: "error",
	})
	if err != nil {
		return err
	}

	ComposeCli, err = compose.NewComposeService(dockerCli)
	if err != nil {
		return err
	}

	return nil
}
