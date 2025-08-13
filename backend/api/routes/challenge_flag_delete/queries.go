package challenge_flag_delete

import (
	"context"
	"trxd/db"
)

func DeleteFlag(ctx context.Context, challengeID int32, flag string) error {
	err := db.Sql.DeleteFlag(ctx, db.DeleteFlagParams{
		ChallID: challengeID,
		Flag:    flag,
	})
	if err != nil {
		return err
	}

	return nil
}
