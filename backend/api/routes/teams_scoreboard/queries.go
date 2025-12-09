package teams_scoreboard

import (
	"context"
	"encoding/json"
	"trxd/db"
)

type TeamData struct {
	ID      int32           `json:"id"`
	Name    string          `json:"name"`
	Score   int32           `json:"score"`
	Country string          `json:"country"`
	Badges  json.RawMessage `json:"badges"`
}

func GetTeamScoreboard(ctx context.Context) ([]*TeamData, error) {
	teams, err := db.Sql.GetTeamsScoreboard(ctx)
	if err != nil {
		return nil, err
	}

	var teamsData []*TeamData
	for _, team := range teams {
		teamData := &TeamData{
			ID:    team.ID,
			Name:  team.Name,
			Score: team.Score,
		}
		if team.Country.Valid {
			teamData.Country = team.Country.String
		}

		if js, ok := team.Badges.([]byte); ok {
			teamData.Badges = js
		}

		teamsData = append(teamsData, teamData)
	}

	return teamsData, nil
}
