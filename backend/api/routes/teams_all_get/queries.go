package teams_all_get

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

type TeamData struct {
	ID      int32                       `json:"id"`
	Name    string                      `json:"name"`
	Score   int32                       `json:"score"`
	Country string                      `json:"country"`
	Image   string                      `json:"image,omitempty"`
	Badges  []sqlc.GetBadgesFromTeamRow `json:"badges,omitempty"`
}

func GetTeams(ctx context.Context) ([]*TeamData, error) {
	teamPreviews, err := db.Sql.GetTeamsPreview(ctx)
	if err != nil {
		return nil, err
	}

	var teamsData []*TeamData
	for _, team := range teamPreviews {
		teamData := &TeamData{
			ID:    team.ID,
			Name:  team.Name,
			Score: team.Score,
		}
		if team.Country.Valid {
			teamData.Country = team.Country.String
		}
		if team.Image.Valid {
			teamData.Image = team.Image.String
		}

		badges, err := db.GetBadgesFromTeam(ctx, team.ID)
		if err != nil {
			return nil, err
		}
		teamData.Badges = badges

		teamsData = append(teamsData, teamData)
	}

	return teamsData, nil
}
