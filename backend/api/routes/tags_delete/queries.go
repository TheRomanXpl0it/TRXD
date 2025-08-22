package tags_delete

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func DeleteTag(ctx context.Context, challID int32, name string) error {
	err := db.Sql.DeleteTag(ctx, sqlc.DeleteTagParams{
		ChallID: challID,
		Name:    name,
	})
	if err != nil {
		return err
	}

	return nil
}
