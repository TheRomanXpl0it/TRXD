package db

import (
	"context"
	"database/sql"
	"trxd/db/sqlc"
)

func GetTeamByID(ctx context.Context, teamID int32) (*sqlc.Team, error) {
	team, err := Sql.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func GetTeamFromUser(ctx context.Context, userID int32) (*sqlc.Team, error) {
	team, err := Sql.GetTeamFromUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func GetTeamByName(ctx context.Context, name string) (*sqlc.Team, error) {
	team, err := Sql.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func GetTotalTeams(ctx context.Context) (int64, error) {
	total, err := Sql.GetTotalTeams(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return total, nil
}
