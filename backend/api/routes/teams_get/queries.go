package teams_get

import (
	"context"
	"database/sql"
	"trxd/db"
)

type TeamData struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Score       int32  `json:"score"`
	Nationality string `json:"nationality"`
	Image       string `json:"image,omitempty"`
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

	return &teamData, nil
}

func GetTeams(ctx context.Context, admin bool) ([]*TeamData, error) {
	teamIDs, err := db.Sql.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	var teams []*TeamData
	for _, teamID := range teamIDs {
		team, err := GetTeam(ctx, teamID, admin)
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
