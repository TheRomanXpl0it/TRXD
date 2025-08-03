package db

import (
	"context"
	"database/sql"

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
