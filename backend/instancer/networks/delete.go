package networks

import (
	"context"
	"strings"
	"trxd/db"
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

	proxyID, err := containers.FetchProxyID(ctx)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "container not found") {
			return err
		}
		proxyID = ""
	}

	if proxyID != "" {
		err = containers.Cli.NetworkDisconnect(ctx, summary[0].ID, proxyID, true)
		if err != nil {
			not_found := strings.Contains(err.Error(), "not found")
			if not_found {
				err := db.StorageDelete(ctx, "proxy-id")
				if err != nil {
					return err
				}
			}
			if !strings.Contains(err.Error(), "is not connected") &&
				!not_found { // endpoint NGINX_ID not found
				return err
			}
		}
	}

	err = containers.Cli.NetworkRemove(ctx, summary[0].ID)
	if err != nil {
		return err
	}

	return nil
}
