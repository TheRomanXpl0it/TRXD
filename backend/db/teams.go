package db

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func AddTeamMember(teamID int32, userID int32) error {
	params := AddTeamMemberParams{
		TeamID: sql.NullInt32{Int32: teamID, Valid: true},
		ID:     userID,
	}
	return queries.AddTeamMember(context.Background(), params)
}

func RegisterTeam(name, password string, user *User) (*Team, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	teamParams := RegisterTeamParams{
		Name:         name,
		PasswordHash: passwordHash,
	}
	team, err := queries.RegisterTeam(context.Background(), teamParams)
	if err != nil {
		return nil, err
	}

	err = AddTeamMember(team.ID, user.ID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}
