package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterTeam(ctx context.Context, name, password string, userID int32) (*Team, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = queries.RegisterTeam(ctx, RegisterTeamParams{
		ID:           userID,
		Name:         name,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	team, err := queries.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			//? this will happen in a race condition on the check if the user is already in a team
			return nil, fmt.Errorf("team %s not found (probably user already in a team)", name)
		}
		return nil, err
	}

	return &team, nil
}

func authTeam(ctx context.Context, qtx *Queries, name, password string) (*Team, error) {
	team, err := qtx.GetTeamByName(ctx, name)
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
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	team, err := authTeam(ctx, qtx, name, password)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	err = qtx.AddTeamMember(ctx, AddTeamMemberParams{
		TeamID: sql.NullInt32{Int32: team.ID, Valid: true},
		ID:     userID,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return team, nil
}

func GetTeamFromUser(ctx context.Context, userID int32) (*Team, error) {
	team, err := queries.GetTeamFromUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}
