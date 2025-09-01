package instancer

import (
	"context"
	"database/sql"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/tde-nico/log"
)

func DeleteInstance(ctx context.Context, tid int32, challID int32, dockerID sql.NullString) error {
	log.Info("Deleting instance:", "team", tid, "challenge", challID)

	if dockerID.Valid {
		var err error
		if len(dockerID.String) == 64 {
			err = KillContainer(ctx, dockerID.String)
		} else {
			err = KillCompose(ctx, dockerID.String)
		}
		if err != nil {
			return err
		}
	}

	err := dbDeleteInstance(ctx, tid, challID)
	if err != nil {
		return err
	}

	return nil
}

func KillContainer(ctx context.Context, id string) error {
	if cli == nil {
		return nil
	}

	err := cli.ContainerRemove(ctx, id, container.RemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func KillCompose(ctx context.Context, name string) error {
	if composeCli == nil {
		return nil
	}

	err := composeCli.Down(ctx, name, api.DownOptions{})
	if err != nil {
		return err
	}

	return nil
}
