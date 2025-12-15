package instancer

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/instancer/composes"
	"trxd/instancer/containers"
	"trxd/instancer/networks"

	"github.com/tde-nico/log"
)

func DeleteInstance(ctx context.Context, tid int32, challID int32, dockerID sql.NullString) error {
	log.Info("Deleting instance:", "challenge", challID, "team", tid)

	if dockerID.Valid {
		var err error
		if len(dockerID.String) == 64 {
			err = containers.KillContainer(ctx, dockerID.String)
		} else {
			err = composes.KillCompose(ctx, dockerID.String)
		}
		if err != nil {
			return err
		}
	}

	err := networks.NetworkDelete(ctx, fmt.Sprintf("net_%d_%d", challID, tid))
	if err != nil {
		return err
	}

	err = dbDeleteInstance(ctx, tid, challID)
	if err != nil {
		return err
	}

	return nil
}
