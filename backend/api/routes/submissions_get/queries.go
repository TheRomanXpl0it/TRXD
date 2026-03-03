package submissions_get

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func GetSubmissions(ctx context.Context, offset int32, limit int32) (int64, []sqlc.GetSubmissionsRow, error) {
	total, err := db.Sql.GetTotalSubmissions(ctx)
	if err != nil {
		return 0, nil, err
	}

	submissions, err := db.Sql.GetSubmissions(ctx, sqlc.GetSubmissionsParams{
		Offset: offset,
		Limit:  sql.NullInt32{Int32: limit, Valid: limit != 0},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return total, []sqlc.GetSubmissionsRow{}, nil
		}
		return 0, nil, err
	}

	// TODO: if user-mode, filter out user-id and user-name

	return total, submissions, err
}
