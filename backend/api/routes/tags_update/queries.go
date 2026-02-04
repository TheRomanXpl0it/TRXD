package tags_update

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateTag(ctx context.Context, challID int32, oldName string, newName string) error {
	err := db.Sql.UpdateTag(ctx, sqlc.UpdateTagParams{
		ChallID: challID,
		OldName: oldName,
		NewName: newName,
	})
	if err != nil {
		return err
	}

	return nil
}
