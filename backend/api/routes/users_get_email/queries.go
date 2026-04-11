package users_get_email

import (
	"context"
	"database/sql"
	"trxd/api/routes/users_get"
	"trxd/db"
)

func GetUserByEmail(ctx context.Context, email string) (*users_get.UserData, error) {
	id, err := db.Sql.GetUserIDByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	data, err := users_get.GetUser(ctx, id, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}
