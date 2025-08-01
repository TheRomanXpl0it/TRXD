package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

func CreateChallenge(ctx context.Context, name, category, description string,
	challType DeployType, maxPoints int32, scoreType ScoreType) (*Challenge, error) {
	id, err := queries.CreateChallenge(ctx, CreateChallengeParams{
		Name:        name,
		Category:    category,
		Description: description,
		Type:        challType,
		MaxPoints:   maxPoints,
		ScoreType:   scoreType,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &Challenge{
		ID:          id,
		Name:        name,
		Category:    category,
		Description: description,
		Type:        challType,
		MaxPoints:   maxPoints,
		ScoreType:   scoreType,
	}, nil
}

func GetChallengeByID(ctx context.Context, challengeID int32) (*Challenge, error) {
	challenge, err := queries.GetChallengeByID(ctx, challengeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &challenge, nil
}

func CreateFlag(ctx context.Context, challengeID int32, flag string, regex bool) (*Flag, error) {
	err := queries.CreateFlag(ctx, CreateFlagParams{
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

	return &Flag{
		ChallID: challengeID,
		Flag:    flag,
		Regex:   regex,
	}, nil
}

func DeleteChallenge(ctx context.Context, challengeID int32) error {
	err := queries.DeleteChallenge(ctx, challengeID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFlag(ctx context.Context, challengeID int32, flag string) error {
	err := queries.DeleteFlag(ctx, DeleteFlagParams{
		ChallID: challengeID,
		Flag:    flag,
	})
	if err != nil {
		return err
	}

	return nil
}
