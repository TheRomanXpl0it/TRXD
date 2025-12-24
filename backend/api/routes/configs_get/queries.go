package configs_get

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func GetConfigs(ctx context.Context) ([]sqlc.Config, error) {
	configs, err := db.Sql.GetConfigs(ctx)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
