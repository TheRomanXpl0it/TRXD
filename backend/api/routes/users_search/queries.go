package users_search

import (
	"context"
	"database/sql"
	"trxd/api/routes/users_get"
	"trxd/db"
)

func GetUserByEmail(ctx context.Context, email string) (*users_get.UserData, error) {
	id, err := db.Sql.GetUserIDByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	data, err := users_get.GetUser(ctx, id, true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetUserByName(ctx context.Context, name string, uid interface{}, allData bool) (*users_get.UserData, error) {
	id, err := db.Sql.GetUserIDByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if uidInt32, ok := uid.(int32); ok {
		allData = allData || uidInt32 == id
	}

	data, err := users_get.GetUser(ctx, id, allData)
	if err != nil {
		return nil, err
	}

	return data, nil
}
