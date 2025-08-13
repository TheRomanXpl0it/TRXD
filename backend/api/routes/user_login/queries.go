package user_login

import (
	"context"
	"database/sql"
	"trxd/db"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(ctx context.Context, email, password string) (*db.User, error) {
	user, err := db.Sql.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
