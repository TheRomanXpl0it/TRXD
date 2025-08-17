package users_get

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
)

type UserData struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Score   int32  `json:"score"`
	Country string `json:"country"`
	Image   string `json:"image"`
}

func GetUsers(ctx context.Context, admin bool) ([]*UserData, error) {
	userPreviews, err := db.Sql.GetUsersPreview(ctx)
	if err != nil {
		return nil, err
	}

	usersData := make([]*UserData, 0)
	for _, user := range userPreviews {
		if !admin && utils.In(user.Role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin}) {
			continue
		}

		userData := &UserData{
			ID:    user.ID,
			Name:  user.Name,
			Score: user.Score,
		}
		if admin {
			userData.Email = user.Email
			userData.Role = string(user.Role)
		}
		if user.Country.Valid {
			userData.Country = user.Country.String
		}
		if user.Image.Valid {
			userData.Image = user.Image.String
		}

		usersData = append(usersData, userData)
	}

	return usersData, nil
}
