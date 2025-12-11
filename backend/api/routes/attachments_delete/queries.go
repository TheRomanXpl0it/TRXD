package attachments_delete

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func DeleteAttachment(ctx context.Context, challID int32, name string) error {
	err := db.Sql.DeleteAttachment(ctx, sqlc.DeleteAttachmentParams{
		ChallID: challID,
		Name:    name,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetAttachmentHash(ctx context.Context, challID int32, name string) (string, error) {
	hash, err := db.Sql.GetAttachmentHash(ctx, sqlc.GetAttachmentHashParams{
		ChallID: challID,
		Name:    name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return hash, nil
}
