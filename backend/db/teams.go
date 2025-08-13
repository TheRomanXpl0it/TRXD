package db

import (
	"context"
	"database/sql"
)

func GetTeamFromUser(ctx context.Context, userID int32) (*Team, error) {
	team, err := Sql.GetTeamFromUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func GetTeamByName(ctx context.Context, name string) (*Team, error) {
	team, err := Sql.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}
