package users_get

import (
	"context"
	"trxd/api/routes/user_get"
	"trxd/db"
)

func GetUsers(ctx context.Context, admin bool) ([]*user_get.UserData, error) {
	userIDs, err := db.Sql.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersData := make([]*user_get.UserData, 0)
	for _, userID := range userIDs {
		userData, err := user_get.GetUser(ctx, userID, admin, true)
		if err != nil {
			return nil, err
		}
		if userData == nil {
			continue
		}
		usersData = append(usersData, userData)
	}

	return usersData, nil
}
