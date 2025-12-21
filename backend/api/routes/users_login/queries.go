package users_login

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/crypto_utils"
)

func LoginUser(ctx context.Context, email, password string) (*sqlc.User, error) {
	user, err := db.Sql.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	valid, err := crypto_utils.Verify(password, user.PasswordSalt, user.PasswordHash)
	if err != nil || !valid {
		return nil, err
	}

	return &user, nil
}
