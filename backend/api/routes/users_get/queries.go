package users_get

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
)

type UserData struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	Score       int32  `json:"score"`
	Nationality string `json:"nationality"`
	Image       string `json:"image"`
}

func GetUser(ctx context.Context, id int32, admin bool) (*UserData, error) {
	data := UserData{}

	user, err := db.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	if !admin && utils.In(user.Role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin}) {
		return nil, nil
	}

	data.ID = user.ID
	data.Name = user.Name
	if admin {
		data.Email = user.Email
		data.Role = string(user.Role)
	}
	data.Score = user.Score
	if user.Nationality.Valid {
		data.Nationality = user.Nationality.String
	}
	if user.Image.Valid {
		data.Image = user.Image.String
	}

	return &data, nil
}

func GetUsers(ctx context.Context, admin bool) ([]*UserData, error) {
	userIDs, err := db.Sql.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersData := make([]*UserData, 0)
	for _, userID := range userIDs {
		userData, err := GetUser(ctx, userID, admin)
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
