package teams_register

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterTeam(ctx context.Context, name, password string, userID int32) (*sqlc.Team, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = db.Sql.RegisterTeam(ctx, sqlc.RegisterTeamParams{
		ID:           userID,
		Name:         name,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	team, err := db.Sql.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("[race condition] team %s not found", name)
		}
		return nil, err
	}

	return &team, nil
}
