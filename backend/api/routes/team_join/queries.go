package team_join

import (
	"context"
	"database/sql"
	"trxd/db"

	"golang.org/x/crypto/bcrypt"
)

func authTeam(ctx context.Context, name, password string) (*db.Team, error) {
	team, err := db.GetTeamByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(team.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, nil
		}
		return nil, err
	}

	return team, nil
}

func JoinTeam(ctx context.Context, name, password string, userID int32) (*db.Team, error) {
	team, err := authTeam(ctx, name, password)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	err = db.Sql.AddTeamMember(ctx, db.AddTeamMemberParams{
		TeamID: sql.NullInt32{Int32: team.ID, Valid: true},
		ID:     userID,
	})
	if err != nil {
		return nil, err
	}

	return team, nil
}
