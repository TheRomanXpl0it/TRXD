package containers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func FetchNginxID(ctx context.Context) (string, error) {
	args := filters.NewArgs()
	// args.Add("name", "trxd-nginx-1") // TODO: make dynamic
	args.Add("name", "trxd-test-pipeline-nginx-1") // TODO: make dynamic
	summary, err := Cli.ContainerList(ctx, container.ListOptions{
		Filters: args,
	})
	if err != nil {
		return "", err
	}
	if len(summary) != 1 {
		return "", fmt.Errorf("expected 1 network, got %d", len(summary))
	}

	return summary[0].ID, nil
}
