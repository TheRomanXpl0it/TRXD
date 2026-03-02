package networks

import (
	"context"
	"errors"
	"trxd/instancer/containers"

	"github.com/docker/docker/api/types/network"
)

func CreateNetwork(ctx context.Context, name string, disableICC bool) (string, error) {
	if containers.Cli == nil {
		return "", nil
	}

	summary, err := FetchNetwork(ctx, name)
	if err != nil {
		return "", err
	}

	var netID string
	if len(summary) > 1 {
		return "", errors.New("multiple networks with the same name") // TODO: test
	} else if len(summary) == 1 {
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

	return netID, nil
}
