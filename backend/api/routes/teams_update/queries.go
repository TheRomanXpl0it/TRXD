package teams_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateTeam(ctx context.Context, tx *sql.Tx, teamID int32, name string, country *string, image *string, bio *string) error {
	params := sqlc.UpdateTeamParams{
		ID:   teamID,
		Name: sql.NullString{String: name, Valid: name != ""},
	}
	if country != nil {
		params.Country = sql.NullString{String: *country, Valid: true}
	}
	if image != nil {
		params.Image = sql.NullString{String: *image, Valid: true}
	}
	if bio != nil {
		params.Bio = sql.NullString{String: *bio, Valid: true}
	}

	sqlTx := db.Sql.WithTx(tx)
	err := sqlTx.UpdateTeam(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
