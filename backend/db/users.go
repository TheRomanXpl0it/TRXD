package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, name, email, password string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	params := RegisterUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         UserRolePlayer,
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

	err = bcrypt.CompareHashAndPassword(user.PasswordHash.([]byte), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
