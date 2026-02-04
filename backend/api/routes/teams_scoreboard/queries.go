package teams_scoreboard

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

func GetTeamScoreboard(ctx context.Context, offset int32, limit int32) (int64, []TeamData, error) {
	total, err := db.GetTotalTeams(ctx)
	if err != nil {
		return 0, nil, err
	}

	teams, err := db.Sql.GetTeamsScoreboard(ctx, sqlc.GetTeamsScoreboardParams{
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
			ID:     team.ID,
			Name:   team.Name,
			Score:  team.Score,
			Badges: team.Badges,
		}
		if team.Country.Valid {
			teamData.Country = team.Country.String
		}

		teamsData = append(teamsData, teamData)
	}

	return total, teamsData, nil
}
