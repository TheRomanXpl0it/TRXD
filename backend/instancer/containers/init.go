package containers

import "github.com/docker/docker/client"

var Cli *client.Client

func InitCli() error {
	var err error

	Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	return nil
}

func CloseCli() {
	if Cli != nil {
		Cli.Close()
	}
}
