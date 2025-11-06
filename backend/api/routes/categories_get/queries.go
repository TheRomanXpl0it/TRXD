package categories_get

import (
	"context"
	"trxd/db"
)

type Category struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func GetCategories(ctx context.Context, author bool) ([]Category, error) {
	categories, err := db.Sql.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	cats := make([]Category, 0, len(categories))
	for _, cat := range categories {
		if !author && cat.VisibleChalls == 0 {
			continue
		}

		cats = append(cats, Category{
			Name: cat.Name,
			Icon: cat.Icon,
		})
	}

	return cats, nil
}
