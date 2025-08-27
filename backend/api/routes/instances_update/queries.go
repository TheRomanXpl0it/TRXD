package instances_update

import (
	"context"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateInstance(ctx context.Context, tid int32, challID int32, expiresAt time.Time) error {
	err := db.Sql.UpdateInstance(ctx, sqlc.UpdateInstanceParams{
		TeamID:    tid,
		ChallID:   challID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	return nil
}
