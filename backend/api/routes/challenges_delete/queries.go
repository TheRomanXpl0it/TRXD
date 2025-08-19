package challenges_delete

import (
	"context"
	"trxd/db"
)

func DeleteChallenge(ctx context.Context, challengeID int32) error {
	err := db.Sql.DeleteChallenge(ctx, challengeID)
	if err != nil {
		return err
	}

	return nil
}
