package db

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/utils"

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
		PasswordHash: string(passwordHash),
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

	err = bcrypt.CompareHashAndPassword([]byte(team.PasswordHash), []byte(password))
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

func GetTeamByName(ctx context.Context, name string) (*Team, error) {
	team, err := queries.GetTeamByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

func UpdateTeam(ctx context.Context, teamID int32, nationality, image, bio string) error {
	err := queries.UpdateTeam(ctx, UpdateTeamParams{
		ID:          teamID,
		Nationality: sql.NullString{String: nationality, Valid: nationality != ""},
		Image:       sql.NullString{String: image, Valid: image != ""},
		Bio:         sql.NullString{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return err
	}

	return nil
}

func ResetTeamPassword(ctx context.Context, teamID int32, newPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = queries.ResetTeamPassword(ctx, ResetTeamPasswordParams{
		ID:           teamID,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return err
	}

	return nil
}

type TeamData struct {
	ID          int32               `json:"id"`
	Name        string              `json:"name"`
	Score       int32               `json:"score"`
	Nationality string              `json:"nationality"`
	Image       string              `json:"image,omitempty"`
	Bio         string              `json:"bio,omitempty"`
	Members     []GetTeamMembersRow `json:"members,omitempty"`
	Solves      []GetTeamSolvesRow  `json:"solves,omitempty"`
	// TODO: add badges
}

func GetTeam(ctx context.Context, teamID int32, admin bool, minimal bool) (*TeamData, error) {
	teamData := TeamData{}

	team, err := queries.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	teamData.ID = team.ID
	teamData.Name = team.Name
	teamData.Score = team.Score
	if team.Nationality.Valid {
		teamData.Nationality = team.Nationality.String
	}
	if team.Image.Valid {
		teamData.Image = team.Image.String
	}

	if minimal {
		return &teamData, nil
	}

	if team.Bio.Valid {
		teamData.Bio = team.Bio.String
	}

	members, err := queries.GetTeamMembers(ctx, sql.NullInt32{Int32: teamID, Valid: true})
	if err != nil {
		return nil, err
	}
	teamData.Members = make([]GetTeamMembersRow, 0)
	for _, member := range members {
		if !admin && utils.In(member.Role, []UserRole{UserRoleAuthor, UserRoleAdmin}) {
			continue
		}
		teamData.Members = append(teamData.Members, member)
	}

	solves, err := queries.GetTeamSolves(ctx, teamID)
	if err != nil {
		return nil, err
	}
	teamData.Solves = solves

	return &teamData, nil
}

func GetTeams(ctx context.Context, admin bool) ([]*TeamData, error) {
	teamIDs, err := queries.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	var teams []*TeamData
	for _, teamID := range teamIDs {
		team, err := GetTeam(ctx, teamID, admin, true)
		if err != nil {
			return nil, err
		}
		if team == nil {
			continue
		}
		teams = append(teams, team)
	}

	return teams, nil
}
