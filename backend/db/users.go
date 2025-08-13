package db

import (
	"context"
	"database/sql"
)

func GetUserByID(ctx context.Context, userID int32) (*User, error) {
	user, err := Sql.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
