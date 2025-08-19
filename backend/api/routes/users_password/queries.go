package users_password

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"

	"golang.org/x/crypto/bcrypt"
)

func ResetUserPassword(ctx context.Context, userID int32, newPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.Sql.ResetUserPassword(ctx, sqlc.ResetUserPasswordParams{
		ID:           userID,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return err
	}

	return nil
}
