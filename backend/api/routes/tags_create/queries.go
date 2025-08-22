package tags_create

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func CreateTag(ctx context.Context, challID int32, name string) error {
	err := db.Sql.CreateTag(ctx, sqlc.CreateTagParams{
		ChallID: challID,
		Name:    name,
	})
	if err != nil {
		return err
	}

	return nil
}
