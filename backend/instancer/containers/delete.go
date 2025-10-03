package containers

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func KillContainer(ctx context.Context, id string) error {
	if Cli == nil {
		return nil
	}

	err := Cli.ContainerRemove(ctx, id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		if !strings.Contains(err.Error(), "No such container") {
			return err
		}
	}

	return nil
}
