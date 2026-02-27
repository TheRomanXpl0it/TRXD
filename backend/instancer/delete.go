package instancer

import (
	"context"
	"database/sql"
	"trxd/instancer/composes"
	"trxd/instancer/containers"

	"trxd/utils/log"
)

func DeleteInstance(ctx context.Context, tid int32, challID int32, dockerID sql.NullString) error {
	log.Info("Deleting instance:", "chall", challID, "team", tid)

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

	err := dbDeleteInstance(ctx, tid, challID)
	if err != nil {
		return err
	}

	return nil
}
