package team_get

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
)

type TeamData struct {
	ID          int32                    `json:"id"`
	Name        string                   `json:"name"`
	Score       int32                    `json:"score"`
	Nationality string                   `json:"nationality"`
	Image       string                   `json:"image,omitempty"`
	Bio         string                   `json:"bio,omitempty"`
	Members     []sqlc.GetTeamMembersRow `json:"members,omitempty"`
	Solves      []sqlc.GetTeamSolvesRow  `json:"solves,omitempty"`
	// TODO: add badges
}

func GetTeam(ctx context.Context, teamID int32, admin bool) (*TeamData, error) {
	teamData := TeamData{}

	team, err := db.GetTeamByID(ctx, teamID)
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

	if team.Bio.Valid {
		teamData.Bio = team.Bio.String
	}

	members, err := db.Sql.GetTeamMembers(ctx, sql.NullInt32{Int32: teamID, Valid: true})
	if err != nil {
		return nil, err
	}
	teamData.Members = make([]sqlc.GetTeamMembersRow, 0)
	for _, member := range members {
		if !admin && utils.In(member.Role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin}) {
			continue
		}
		teamData.Members = append(teamData.Members, member)
	}

	solves, err := db.Sql.GetTeamSolves(ctx, teamID)
	if err != nil {
		return nil, err
	}
	teamData.Solves = solves

	return &teamData, nil
}
