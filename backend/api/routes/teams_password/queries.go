package teams_password

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/crypto_utils"
)

func ResetTeamPassword(ctx context.Context, teamID int32, newPassword string) error {
	hash, salt, err := crypto_utils.Hash(newPassword)
	if err != nil {
		return err
	}

	err = db.Sql.ResetTeamPassword(ctx, sqlc.ResetTeamPasswordParams{
		ID:           teamID,
		PasswordHash: hash,
		PasswordSalt: salt,
	})
	if err != nil {
		return err
	}

	return nil
}
