package users_role

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func ChangeUserRole(ctx context.Context, uid int32, newRole sqlc.UserRole) error {
	err := db.Sql.ChangeUserRole(ctx, sqlc.ChangeUserRoleParams{
		ID:   uid,
		Role: newRole,
	})
	if err != nil {
		return err
	}

	return nil
}
