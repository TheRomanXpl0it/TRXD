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
	Role    string          `json:"role,omitempty"`
	Score   int32           `json:"score"`
	Country string          `json:"country"`
	Badges  json.RawMessage `json:"badges"`
}

func GetTeams(ctx context.Context, offset int32, limit int32) (int64, []TeamData, error) {
	modeStr, err := db.GetConfig(ctx, "user-mode")
	if err != nil {
		return 0, nil, err
	}
	userMode := modeStr == "true"

	total, err := db.GetTotalTeams(ctx)
	if err != nil {
		return 0, nil, err
	}

	teams, err := db.Sql.GetTeamsPreview(ctx, sqlc.GetTeamsPreviewParams{
		Offset: offset,
		Limit:  sql.NullInt32{Int32: limit, Valid: limit != 0},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return total, []TeamData{}, nil
		}
		return 0, nil, err
	}

	var teamsData []TeamData
	for _, team := range teams {
		teamData := TeamData{
			ID:    team.ID,
			Name:  team.Name,
			Score: team.Score,
		}
		if userMode {
			teamData.Role = string(team.UserRole)
		}
		if team.Country.Valid {
			teamData.Country = team.Country.String
		}

		if js, ok := team.Badges.([]byte); ok {
			teamData.Badges = js
		}

		teamsData = append(teamsData, teamData)
	}

	return total, teamsData, nil
}
