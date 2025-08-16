package db

import (
	"context"
	"database/sql"
	"trxd/db/sqlc"
)

func GetChallengeByID(ctx context.Context, challengeID int32) (*sqlc.Challenge, error) {
	challenge, err := Sql.GetChallengeByID(ctx, challengeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &challenge, nil
}

func GetTagsByChallenge(ctx context.Context, challengeID int32) ([]string, error) {
	tags, err := Sql.GetTagsByChallenge(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func IsChallengeSolved(ctx context.Context, id int32, uid int32) (bool, error) {
	solved, err := Sql.IsChallengeSolved(ctx, sqlc.IsChallengeSolvedParams{
		ChallID: id,
		ID:      uid,
	})
	if err != nil {
		return false, err
	}

	return solved, nil
}
