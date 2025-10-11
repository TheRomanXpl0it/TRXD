package teams_register

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/crypto_utils"

	"github.com/lib/pq"
)

func RegisterTeam(ctx context.Context, tx *sql.Tx, name, password string, userID int32) (*sqlc.Team, error) {
	hash, salt, err := crypto_utils.Hash(password)
	if err != nil {
		return nil, err
	}

	sqlTx := db.Sql.WithTx(tx)
	err = sqlTx.RegisterTeam(ctx, sqlc.RegisterTeamParams{
		ID:           userID,
		Name:         name,
		PasswordHash: hash,
		PasswordSalt: salt,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	team, err := sqlTx.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("[race condition] team %s not found", name)
		}
		return nil, err
	}

	return &team, nil
}
