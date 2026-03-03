package submissions_delete

import (
	"context"
	"trxd/db"
)

func DeleteSubmission(ctx context.Context, subID int32) error {
	err := db.Sql.DeleteSubmission(ctx, subID)
	if err != nil {
		return err
	}

	return nil
}
