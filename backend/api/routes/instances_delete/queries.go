package instances_delete

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

type Chall struct {
	Info         *sqlc.Challenge
	DockerConfig *sqlc.GetDockerConfigsByIDRow
}

// TODO: move these

func GetChallenge(ctx context.Context, challID int32) (*Chall, error) {
	info := &Chall{}

	chall, err := db.GetChallengeByID(ctx, challID)
	if err != nil {
		return nil, err
	}
	if chall == nil {
		return nil, nil
	}
	info.Info = chall

	dockerConfig, err := db.GetDockerConfigsByID(ctx, challID)
	if err != nil {
		return nil, err
	}
	if dockerConfig == nil {
		return info, nil
	}
	info.DockerConfig = dockerConfig

	return info, nil
}

func GetInstance(ctx context.Context, challID, teamID int32) (*sqlc.Instance, error) {
	instance, err := db.Sql.GetInstance(ctx, sqlc.GetInstanceParams{
		ChallID: challID,
		TeamID:  teamID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &instance, nil
}
