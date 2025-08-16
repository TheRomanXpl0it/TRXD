package challenge_flag_create

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

func CreateFlag(ctx context.Context, challengeID int32, flag string, regex bool) (*sqlc.Flag, error) {
	err := db.Sql.CreateFlag(ctx, sqlc.CreateFlagParams{
		Flag:    flag,
		ChallID: challengeID,
		Regex:   regex,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &sqlc.Flag{
		ChallID: challengeID,
		Flag:    flag,
		Regex:   regex,
	}, nil
}
