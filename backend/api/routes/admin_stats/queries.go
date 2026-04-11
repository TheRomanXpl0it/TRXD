package admin_stats

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func GetAdminStats(ctx context.Context) (sqlc.GetAdminStatsRow, error) {
	stats, err := db.Sql.GetAdminStats(ctx)
	if err != nil {
		return sqlc.GetAdminStatsRow{}, err
	}

	return stats, nil
}
