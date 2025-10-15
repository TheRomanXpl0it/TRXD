package networks

import (
	"context"
	"strings"
	"trxd/instancer/containers"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

func CreateNetwork(ctx context.Context, name string, disableICC bool) (string, error) {
	if containers.Cli == nil {
		return "", nil
	}

	args := filters.NewArgs()
	args.Add("name", name)
	summary, err := containers.Cli.NetworkList(ctx, network.ListOptions{
		Filters: args,
	})
	if err != nil {
		return "", err
	}

	var netID string
	if len(summary) == 1 {
		netID = summary[0].ID
	} else {
		options := make(map[string]string)
		if disableICC {
			options["com.docker.network.bridge.enable_icc"] = "false"
		}

		net, err := containers.Cli.NetworkCreate(ctx, name, network.CreateOptions{
			Internal: !disableICC,
			Options:  options,
		})
		if err != nil {
			return "", err
		}
		netID = net.ID
	}

	if disableICC {
		return netID, nil
	}

	nginxID, err := containers.FetchNginxID(ctx)
	if err != nil {
		return "", err
	}

	err = containers.Cli.NetworkConnect(ctx, netID, nginxID, nil)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists in network") {
			return "", err
		}
	}

	return netID, nil
}
