package category_create

import (
	"context"
	"trxd/db"

	"github.com/lib/pq"
)

func CreateCategory(ctx context.Context, name string, icon string) (*db.Category, error) {
	err := db.Sql.CreateCategory(ctx, db.CreateCategoryParams{
		Name: name,
		Icon: icon,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}
	return &db.Category{Name: name, Icon: icon}, nil
}
