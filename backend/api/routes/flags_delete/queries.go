package flags_delete

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func DeleteFlag(ctx context.Context, challengeID int32, flag string) error {
	err := db.Sql.DeleteFlag(ctx, sqlc.DeleteFlagParams{
		ChallID: challengeID,
		Flag:    flag,
	})
	if err != nil {
		return err
	}

	return nil
}
