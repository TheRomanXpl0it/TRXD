package networks

import (
	"context"
	"trxd/instancer/containers"

	"github.com/docker/docker/api/types/network"
)

func CreateNetwork(ctx context.Context, name string) (string, error) {
	if containers.Cli == nil {
		return "", nil
	}

	// TODO: already exists check

	net, err := containers.Cli.NetworkCreate(ctx, name, network.CreateOptions{
		Internal: true,
	})
	if err != nil {
		return "", err
	}

	nginxID, err := containers.FetchNginxID(ctx)
	if err != nil {
		return "", err
	}

	err = containers.Cli.NetworkConnect(ctx, net.ID, nginxID, nil)
	if err != nil {
		return "", err
	}

	return net.ID, nil
}
