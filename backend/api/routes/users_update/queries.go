package users_update

import (
	"context"
	"database/sql"
	"errors"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

func UpdateUser(ctx context.Context, tx *sql.Tx, userID int32, name string, country *string) error {
	params := sqlc.UpdateUserParams{
		ID:   userID,
		Name: sql.NullString{String: name, Valid: name != ""},
	}
	if country != nil {
		params.Country = sql.NullString{String: *country, Valid: true}
	}
	sqlTx := db.Sql.WithTx(tx)
	err := sqlTx.UpdateUser(ctx, params)
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
