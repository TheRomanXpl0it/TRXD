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

func GetDockerConfigsByID(ctx context.Context, challengeID int32) (*sqlc.GetDockerConfigsByIDRow, error) {
	dockerConfig, err := Sql.GetDockerConfigsByID(ctx, challengeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &dockerConfig, nil
}

func GetTagsByChallenge(ctx context.Context, challengeID int32) ([]string, error) {
	tags, err := Sql.GetTagsByChallenge(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
