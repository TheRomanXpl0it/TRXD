package db

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(name, email, password string, role UserRole) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	params := RegisterUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}
	user, err := queries.RegisterUser(context.Background(), params)
	return &user, err
}
