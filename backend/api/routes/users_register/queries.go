package users_register

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"

	"github.com/lib/pq"
)

func DBRegisterUser(ctx context.Context, tx *sql.Tx, name, email, password string, role ...sqlc.UserRole) (*sqlc.User, error) {
	hash, salt, err := crypto_utils.Hash(password)
	if err != nil {
		return nil, err
	}

	if len(role) == 0 {
		role = append(role, sqlc.UserRolePlayer)
	}

	sqlTx := db.Sql.WithTx(tx)
	user, err := sqlTx.RegisterUser(ctx, sqlc.RegisterUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		PasswordSalt: salt,
		Role:         role[0],
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == consts.PGUniqueViolation {
				return nil, nil
			}
		}
		return nil, err
	}

	return &user, nil
}
