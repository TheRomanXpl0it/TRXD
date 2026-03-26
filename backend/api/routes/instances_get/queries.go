package instances_get

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func GetInstances(ctx context.Context) ([]sqlc.GetInstancesRow, error) {
	instances, err := db.Sql.GetInstances(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return []sqlc.GetInstancesRow{}, nil
		}
		return nil, err
	}
	if instances == nil {
		instances = []sqlc.GetInstancesRow{}
	}

	return instances, nil
}
