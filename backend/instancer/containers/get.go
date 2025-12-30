package containers

import (
	"context"
	"fmt"
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
		return "", fmt.Errorf("container not found (%s)", name)
	}

	return summary[0].ID, nil
}

func FetchProxyID(ctx context.Context) (string, error) {
	id, err := db.StorageGet(ctx, "proxy-id")
	if err != nil {
		return "", err
	}
	if id != nil && *id != "" {
		return *id, nil
	}

	name, err := db.GetConfig(ctx, "project-name")
	if err != nil {
		return "", err
	}
	if name == "" {
		name = consts.DefaultConfigs["project-name"].(string)
	}

	containerID, err := FetchContainerByName(ctx, name+"-nginx-1")
	// containerID, err := FetchContainerByName(ctx, name+"-traefik-1") // TODO: mayber put this as const
	if err != nil {
		return "", err
	}

	err = db.StorageSet(ctx, "proxy-id", containerID)
	if err != nil {
		return "", err
	}

	return containerID, nil
}
