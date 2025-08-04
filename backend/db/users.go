package db

import (
	"context"
	"database/sql"
	"time"
	"trxd/utils"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, name, email, password string, role ...UserRole) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if len(role) == 0 {
		role = append(role, UserRolePlayer)
	}

	params := RegisterUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role[0],
	}
	user, err := queries.RegisterUser(ctx, params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &user, nil
}

func LoginUser(ctx context.Context, email, password string) (*User, error) {
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByID(ctx context.Context, userID int32) (*User, error) {
	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func UpdateUser(ctx context.Context, userID int32, name, nationality, image string) error {
	err := queries.UpdateUser(ctx, UpdateUserParams{
		ID:          userID,
		Name:        sql.NullString{String: name, Valid: name != ""},
		Nationality: sql.NullString{String: nationality, Valid: nationality != ""},
		Image:       sql.NullString{String: image, Valid: image != ""},
	})
	if err != nil {
		return err
	}

	return nil
}

func ResetUserPassword(ctx context.Context, userID int32, newPassword string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = queries.ResetUserPassword(ctx, ResetUserPasswordParams{
		ID:           userID,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return err
	}

	return nil
}

type UserData struct {
	ID          int32              `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Role        string             `json:"role"`
	Score       int32              `json:"score"`
	Nationality string             `json:"nationality"`
	Image       string             `json:"image"`
	TeamID      int32              `json:"team_id"`
	JoinedAt    *time.Time         `json:"joined_at,omitempty"`
	Solves      []GetUserSolvesRow `json:"solves,omitempty"`
}

func GetUser(ctx context.Context, id int32, admin bool, minimal bool) (*UserData, error) {
	data := UserData{}

	user, err := queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if !admin && utils.In(user.Role, []UserRole{UserRoleAuthor, UserRoleAdmin}) {
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

	if minimal {
		return &data, nil
	}

	data.JoinedAt = &user.CreatedAt

	data.TeamID = -1
	if user.TeamID.Valid {
		data.TeamID = user.TeamID.Int32
	}

	solves, err := queries.GetUserSolves(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	data.Solves = []GetUserSolvesRow{}
	if solves != nil {
		data.Solves = solves
	}

	return &data, nil
}

func GetUsers(ctx context.Context, admin bool) ([]*UserData, error) {
	userIDs, err := queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersData := make([]*UserData, 0)
	for _, userID := range userIDs {
		userData, err := GetUser(ctx, userID, admin, true)
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
