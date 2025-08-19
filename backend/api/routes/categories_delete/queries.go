package categories_delete

import (
	"context"
	"trxd/db"
)

func DeleteCategory(ctx context.Context, category string) error {
	err := db.Sql.DeleteCategory(ctx, category)
	if err != nil {
		return err
	}

	return nil
}
