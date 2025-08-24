package instances_update

import (
	"context"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
)

type Chall struct {
	Info         *sqlc.Challenge
	DockerConfig *sqlc.GetDockerConfigsByIDRow
}

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

func UpdateInstance(ctx context.Context, tid int32, challID int32, expiresAt time.Time) error {
	err := db.Sql.UpdateInstance(ctx, sqlc.UpdateInstanceParams{
		TeamID:    tid,
		ChallID:   challID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	return nil
}
