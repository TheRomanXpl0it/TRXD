package teams_password

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"

	"golang.org/x/crypto/bcrypt"
)

func ResetTeamPassword(ctx context.Context, teamID int32, newPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.Sql.ResetTeamPassword(ctx, sqlc.ResetTeamPasswordParams{
		ID:           teamID,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return err
	}

	return nil
}
