package user_register

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, name, email, password string, role ...sqlc.UserRole) (*sqlc.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if len(role) == 0 {
		role = append(role, sqlc.UserRolePlayer)
	}

	user, err := db.Sql.RegisterUser(ctx, sqlc.RegisterUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role[0],
	})
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
