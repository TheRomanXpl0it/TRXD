package team_update

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateTeam(ctx context.Context, teamID int32, country, image, bio string) error {
	err := db.Sql.UpdateTeam(ctx, sqlc.UpdateTeamParams{
		ID:      teamID,
		Country: sql.NullString{String: country, Valid: country != ""},
		Image:   sql.NullString{String: image, Valid: image != ""},
		Bio:     sql.NullString{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return err
	}

	return nil
}
