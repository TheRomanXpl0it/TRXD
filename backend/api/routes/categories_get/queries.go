package categories_get

import (
	"context"
	"trxd/db"
)

func GetCategories(ctx context.Context, author bool) ([]string, error) {
	categories, err := db.Sql.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	cats := make([]string, 0, len(categories))
	for _, cat := range categories {
		if !author && cat.VisibleChalls == 0 {
			continue
		}

		cats = append(cats, cat.Name)
	}

	return cats, nil
}
