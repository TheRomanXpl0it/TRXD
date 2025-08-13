package teams_get

import (
	"context"
	"trxd/api/routes/team_get"
	"trxd/db"
)

func GetTeams(ctx context.Context, admin bool) ([]*team_get.TeamData, error) {
	teamIDs, err := db.Sql.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	var teams []*team_get.TeamData
	for _, teamID := range teamIDs {
		team, err := team_get.GetTeam(ctx, teamID, admin, true)
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
