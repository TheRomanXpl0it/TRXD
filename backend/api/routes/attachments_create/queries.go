package attachments_create

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func CreateAttachments(ctx context.Context, challID int32, names []string, hashes []string) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer db.Rollback(tx)

	sqlx := db.Sql.WithTx(tx)
	for i := range len(names) {
		err = sqlx.CreateAttachment(ctx, sqlc.CreateAttachmentParams{
			ChallID: challID,
			Name:    names[i],
			Hash:    hashes[i],
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
