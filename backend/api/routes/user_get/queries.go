package user_get

import (
	"context"
	"database/sql"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
)

type UserData struct {
	ID       int32                   `json:"id"`
	Name     string                  `json:"name"`
	Email    string                  `json:"email"`
	Role     string                  `json:"role"`
	Score    int32                   `json:"score"`
	Country  string                  `json:"country"`
	Image    string                  `json:"image"`
	TeamID   *int32                  `json:"team_id"`
	JoinedAt *time.Time              `json:"joined_at,omitempty"`
	Solves   []sqlc.GetUserSolvesRow `json:"solves,omitempty"`
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
	if user.Country.Valid {
		data.Country = user.Country.String
	}
	if user.Image.Valid {
		data.Image = user.Image.String
	}

	data.JoinedAt = &user.CreatedAt

	data.TeamID = nil
	if user.TeamID.Valid {
		data.TeamID = &user.TeamID.Int32
	}

	solves, err := db.Sql.GetUserSolves(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	data.Solves = []sqlc.GetUserSolvesRow{}
	if solves != nil {
		data.Solves = solves
	}

	return &data, nil
}
