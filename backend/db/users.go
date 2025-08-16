package db

import (
	"context"
	"database/sql"
	"trxd/db/sqlc"
)

func GetUserByID(ctx context.Context, userID int32) (*sqlc.User, error) {
	user, err := Sql.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
