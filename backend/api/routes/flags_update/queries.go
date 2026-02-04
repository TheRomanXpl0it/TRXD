package flags_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"

	"github.com/lib/pq"
)

func UpdateFlag(ctx context.Context, challID int32, flag string, regex *bool, newFlag string) (bool, error) {
	nullBool := sql.NullBool{
		Valid: regex != nil,
	}
	if nullBool.Valid {
		nullBool.Bool = *regex
	}

	err := db.Sql.UpdateFlag(ctx, sqlc.UpdateFlagParams{
		ChallID: challID,
		Flag:    flag,
		Regex:   nullBool,
		NewFlag: sql.NullString{
			String: newFlag,
			Valid:  newFlag != flag && newFlag != "",
		},
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == consts.PGUniqueViolation {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}
