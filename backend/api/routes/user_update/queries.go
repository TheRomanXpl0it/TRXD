package user_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateUser(ctx context.Context, userID int32, name, country, image string) error {
	err := db.Sql.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:      userID,
		Name:    sql.NullString{String: name, Valid: name != ""},
		Country: sql.NullString{String: country, Valid: country != ""},
		Image:   sql.NullString{String: image, Valid: image != ""},
	})
	if err != nil {
		return err
	}

	return nil
}
