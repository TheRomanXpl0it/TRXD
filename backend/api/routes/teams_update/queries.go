package teams_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateTeam(ctx context.Context, tx *sql.Tx, teamID int32, name, country, image, bio string) error {
	sqlTx := db.Sql.WithTx(tx)
	err := sqlTx.UpdateTeam(ctx, sqlc.UpdateTeamParams{
		ID:      teamID,
		Name:    sql.NullString{String: name, Valid: name != ""},
		Country: sql.NullString{String: country, Valid: country != ""},
		Image:   sql.NullString{String: image, Valid: image != ""},
		Bio:     sql.NullString{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return err
	}

	return nil
}
