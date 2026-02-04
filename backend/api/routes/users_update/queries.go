package users_update

import (
	"context"
	"database/sql"
	"errors"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"

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
			if pqErr.Code == consts.PGUniqueViolation {
				return errors.New("[name already taken]")
			}
		}
		return err
	}

	return nil
}
