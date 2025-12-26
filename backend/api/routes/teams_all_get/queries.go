package teams_all_get

import (
	"context"
	"database/sql"
	"encoding/json"
	"trxd/db"
	"trxd/db/sqlc"
)

type TeamData struct {
	ID      int32           `json:"id"`
	Name    string          `json:"name"`
	Score   int32           `json:"score"`
	Country string          `json:"country"`
	Badges  json.RawMessage `json:"badges"`
}

func GetTeams(ctx context.Context, start int32, end int32) ([]*TeamData, error) {
	teams, err := db.Sql.GetTeamsPreview(ctx, sqlc.GetTeamsPreviewParams{
		Offset: start,
		Limit:  sql.NullInt32{Int32: end - start, Valid: end != 0},
	})
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
