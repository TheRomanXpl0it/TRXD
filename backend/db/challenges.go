package db

import (
	"context"
	"database/sql"
)

func GetChallengeByID(ctx context.Context, challengeID int32) (*Challenge, error) {
	challenge, err := Sql.GetChallengeByID(ctx, challengeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &challenge, nil
}
