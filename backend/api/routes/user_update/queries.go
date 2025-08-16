package user_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateUser(ctx context.Context, userID int32, name, nationality, image string) error {
	err := db.Sql.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          userID,
		Name:        sql.NullString{String: name, Valid: name != ""},
		Nationality: sql.NullString{String: nationality, Valid: nationality != ""},
		Image:       sql.NullString{String: image, Valid: image != ""},
	})
	if err != nil {
		return err
	}

	return nil
}
