package db

import (
	"context"
	"database/sql"
	"time"
	"trxd/db/sqlc"
)

type Chall struct {
	Info         *sqlc.Challenge
	DockerConfig *sqlc.GetDockerConfigsByIDRow
}

func GetChallenge(ctx context.Context, challID int32) (*Chall, error) {
	info := &Chall{}

	chall, err := GetChallengeByID(ctx, challID)
	if err != nil {
		return nil, err
	}
	if chall == nil {
		return nil, nil
	}
	info.Info = chall

	dockerConfig, err := GetDockerConfigsByID(ctx, challID)
	if err != nil {
		return nil, err
	}
	if dockerConfig == nil {
		return info, nil
	}
	info.DockerConfig = dockerConfig

	return info, nil
}

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

func GetHiddenAndAttachments(ctx context.Context, challengeID int32) (*sqlc.GetHiddenAndAttachmentsRow, error) {
	res, err := Sql.GetHiddenAndAttachments(ctx, challengeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func GetTotalCategoryChallenges(ctx context.Context) ([]sqlc.GetTotalCategoryChallengesRow, error) {
	start, err := GetConfig(ctx, "start-time")
	if err != nil {
		return nil, err
	}

	if start != "" {
		startTime, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return nil, err
		}
		if time.Now().Before(startTime) {
			return nil, nil
		}
	}

	challenges, err := Sql.GetTotalCategoryChallenges(ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if challenges == nil {
		challenges = make([]sqlc.GetTotalCategoryChallengesRow, 0)
	}

	return challenges, nil
}
