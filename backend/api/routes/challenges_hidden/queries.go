package challenges_hidden

import (
	"context"
	"trxd/db"
)

func ToggleChallengesHidden(ctx context.Context, challIDs []int32) error {
	err := db.Sql.ToggleChallengesHidden(ctx, challIDs)
	if err != nil {
		return err
	}

	return nil
}
