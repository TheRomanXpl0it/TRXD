package db

import (
	"context"

	"github.com/lib/pq"
)

func CreateCategory(ctx context.Context, name string, icon string) (*Category, error) {
	err := queries.CreateCategory(ctx, CreateCategoryParams{
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
	return &Category{Name: name, Icon: icon}, nil
}

func DeleteCategory(ctx context.Context, categoryName string) error {
	err := queries.DeleteCategory(ctx, categoryName)
	if err != nil {
		return err
	}

	return nil
}
