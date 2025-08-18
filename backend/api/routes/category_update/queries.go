package category_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func GetCategory(ctx context.Context, name string) (*sqlc.Category, error) {
	category, err := db.Sql.GetCategory(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func UpdateChallengesCategory(ctx context.Context, oldName string, newName string) error {
	err := db.Sql.UpdateChallengesCategory(ctx, sqlc.UpdateChallengesCategoryParams{
		OldCategory: oldName,
		NewCategory: newName,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateCategoryIcon(ctx context.Context, name string, newIcon string) error {
	err := db.Sql.UpdateCategoryIcon(ctx, sqlc.UpdateCategoryIconParams{
		Name: name,
		Icon: newIcon,
	})
	if err != nil {
		return err
	}

	return nil
}
