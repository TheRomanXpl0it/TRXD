package categories_create

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"

	"github.com/lib/pq"
)

func CreateCategory(ctx context.Context, name string) (*sqlc.Category, error) {
	err := db.Sql.CreateCategory(ctx, name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == consts.PGUniqueViolation {
				return nil, nil
			}
		}
		return nil, err
	}
	return &sqlc.Category{Name: name}, nil
}
