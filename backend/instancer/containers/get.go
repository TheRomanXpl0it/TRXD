package containers

import (
	"context"
	"errors"
	"trxd/db"
	"trxd/utils/consts"

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
		return "", errors.New("nginx container not found")
	}

	return summary[0].ID, nil
}

func FetchNginxID(ctx context.Context) (string, error) {
	// TODO: put a cache here

	name, err := db.GetConfig(ctx, "project-name")
	if err != nil {
		return "", err
	}
	if name == "" {
		name = consts.DefaultConfigs["project-name"].(string) + "-nginx-1"
	}
	return FetchContainerByName(ctx, name)
}
