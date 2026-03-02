package networks

import (
	"context"
	"trxd/instancer/containers"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func FetchNetwork(ctx context.Context, name string) ([]network.Summary, error) {
	args := filters.NewArgs()
	args.Add("name", name)

	summary, err := containers.Cli.NetworkList(ctx, network.ListOptions{
		Filters: args,
	})
	if err != nil {
		return nil, err
	}

	return summary, nil
}
