package users_password

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/crypto_utils"
)

func ResetUserPassword(ctx context.Context, userID int32, newPassword string) error {
	hash, salt, err := crypto_utils.Hash(newPassword)
	if err != nil {
		return err
	}

	err = db.Sql.ResetUserPassword(ctx, sqlc.ResetUserPasswordParams{
		ID:           userID,
		PasswordHash: hash,
		PasswordSalt: salt,
	})
	if err != nil {
		return err
	}

	return nil
}
