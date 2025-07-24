package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func teamMemberAdd(ctx context.Context, teamID int32, userID int32) error {
	params := AddTeamMemberParams{
		TeamID: sql.NullInt32{Int32: teamID, Valid: true},
		ID:     userID,
	}
	return queries.AddTeamMember(ctx, params)
}

func RegisterTeam(ctx context.Context, name, password string, userID int32) (*Team, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	teamParams := RegisterTeamParams{
		Name:         name,
		PasswordHash: passwordHash,
	}
	team, err := queries.RegisterTeam(ctx, teamParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	err = teamMemberAdd(ctx, team.ID, userID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func loginTeam(ctx context.Context, name, password string) (*Team, error) {
	team, err := queries.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(team.PasswordHash.([]byte), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func JoinTeam(ctx context.Context, name, password string, userID int32) (*Team, error) {
	team, err := loginTeam(ctx, name, password)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	err = teamMemberAdd(ctx, team.ID, userID)
	if err != nil {
		return nil, err
	}

	return team, nil
}
