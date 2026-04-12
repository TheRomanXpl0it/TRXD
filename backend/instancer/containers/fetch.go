package containers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func FetchContainerByName(ctx context.Context, name string) (string, error) {
	args := filters.NewArgs()
	args.Add("name", name)
	summary, err := Cli.ContainerList(ctx, container.ListOptions{
		Filters: args,
	})
	if err != nil {
		return "", err
	}
	if len(summary) != 1 {
		return "", fmt.Errorf("container not found (%s)", name)
	}

	return summary[0].ID, nil
}
