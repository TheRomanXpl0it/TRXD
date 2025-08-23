package instances_create

import (
	"context"
	"database/sql"
	"time"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

type Chall struct {
	Info         *sqlc.Challenge
	DockerConfig *sqlc.DockerConfig
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

func CreateInstance(ctx context.Context, teamID, challID int32, expiresAt time.Time, hashDomain bool) (*sqlc.CreateInstanceRow, error) {
	info, err := db.Sql.CreateInstance(ctx, sqlc.CreateInstanceParams{
		TeamID:     teamID,
		ChallID:    challID,
		ExpiresAt:  expiresAt,
		HashDomain: hashDomain,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &info, nil
}
