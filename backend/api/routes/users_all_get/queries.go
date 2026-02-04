package users_all_get

import (
	"context"
	"database/sql"
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
}

func GetUsers(ctx context.Context, isAdmin bool, offset int32, limit int32) (int64, []UserData, error) {
	total, err := db.Sql.GetTotalUsers(ctx, isAdmin)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, nil, err
		}
		total = 0
	}

	userPreviews, err := db.Sql.GetUsers(ctx, sqlc.GetUsersParams{
		IsAdmin: isAdmin,
		Offset:  offset,
		Limit:   sql.NullInt32{Int32: limit, Valid: limit != 0},
	})
	if err != nil {
		if err != sql.ErrNoRows {
			return total, nil, err
		}
	}

	usersData := make([]UserData, 0)
	for _, user := range userPreviews {
		if !isAdmin && utils.In(user.Role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin}) {
			continue
		}

		userData := UserData{
			ID:    user.ID,
			Name:  user.Name,
			Score: user.Score,
		}
		if isAdmin {
			userData.Email = user.Email
			userData.Role = string(user.Role)
		}
		if user.Country.Valid {
			userData.Country = user.Country.String
		}

		usersData = append(usersData, userData)
	}

	return total, usersData, nil
}
