package teams_search

import (
	"context"
	"database/sql"
	"trxd/api/routes/teams_get"
	"trxd/db"
)

func GetTeamByEmail(ctx context.Context, email string) (*teams_get.TeamData, error) {
	idNull, err := db.Sql.GetTeamIDByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !idNull.Valid {
		return nil, nil
	}
	id := idNull.Int32

	data, err := teams_get.GetTeam(ctx, id, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTeamByName(ctx context.Context, name string, tid interface{}, allData bool) (*teams_get.TeamData, error) {
	id, err := db.Sql.GetTeamIDByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if tidInt32, ok := tid.(int32); ok {
		allData = allData || tidInt32 == id
	}

	data, err := teams_get.GetTeam(ctx, id, allData)
	if err != nil {
		return nil, err
	}

	return data, nil
}
