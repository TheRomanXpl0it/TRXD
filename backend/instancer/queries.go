package instancer

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"

	"github.com/lib/pq"
)

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

func dbCreateInstance(ctx context.Context, teamID, challID int32,
	expiresAt time.Time, hashDomain bool) (*sqlc.CreateInstanceRow, error) {

	info, err := db.Sql.CreateInstance(ctx, sqlc.CreateInstanceParams{
		TeamID:     teamID,
		ChallID:    challID,
		ExpiresAt:  expiresAt,
		HashDomain: hashDomain,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == consts.PGUniqueViolation {
				if pqErr.Constraint == "instances_port_key" {
					return nil, errors.New("[port conflict]")
				}
				return nil, nil
			}
		}
		return nil, err
	}

	return &info, nil
}

func dbUpdateInstanceDockerID(ctx context.Context, teamID, challID int32, dockerID string) error {
	err := db.Sql.UpdateInstanceDockerID(ctx, sqlc.UpdateInstanceDockerIDParams{
		TeamID:   teamID,
		ChallID:  challID,
		DockerID: sql.NullString{String: dockerID, Valid: dockerID != ""},
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateInstanceExpire(ctx context.Context, tid int32, challID int32, expiresAt time.Time) error {
	err := db.Sql.UpdateInstanceExpire(ctx, sqlc.UpdateInstanceExpireParams{
		TeamID:    tid,
		ChallID:   challID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func dbDeleteInstance(ctx context.Context, tid int32, challID int32) error {
	err := db.Sql.DeleteInstance(ctx, sqlc.DeleteInstanceParams{
		TeamID:  tid,
		ChallID: challID,
	})
	if err != nil {
		return err
	}

	return nil
}
