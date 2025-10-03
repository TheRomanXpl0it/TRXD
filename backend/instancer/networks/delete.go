package networks

import (
	"context"
	"strings"
	"trxd/instancer/containers"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func NetworkDelete(ctx context.Context, name string) error {
	if containers.Cli == nil {
		return nil
	}

	args := filters.NewArgs()
	args.Add("name", name)
	summary, err := containers.Cli.NetworkList(ctx, network.ListOptions{
		Filters: args,
	})
	if err != nil {
		return err
	}
	if len(summary) != 1 {
		return nil
	}

	nginxID, err := containers.FetchNginxID(ctx)
	if err != nil {
		return err
	}

	err = containers.Cli.NetworkDisconnect(ctx, summary[0].ID, nginxID, true)
	if err != nil {
		if !strings.Contains(err.Error(), "is not connected") {
			return err
		}
	}

	err = containers.Cli.NetworkRemove(ctx, summary[0].ID)
	if err != nil {
		return err
	}

	return nil
}
