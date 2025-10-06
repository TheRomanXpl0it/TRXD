package users_update

import (
	"context"
	"database/sql"
	"errors"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

func UpdateUser(ctx context.Context, userID int32, name, country, image string) error {
	err := db.Sql.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:      userID,
		Name:    sql.NullString{String: name, Valid: name != ""},
		Country: sql.NullString{String: country, Valid: country != ""},
		Image:   sql.NullString{String: image, Valid: image != ""},
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return errors.New("[name already taken]")
			}
		}
		return err
	}

	return nil
}
