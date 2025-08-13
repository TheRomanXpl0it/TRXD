package team_update

import (
	"context"
	"database/sql"
	"trxd/db"
)

func UpdateTeam(ctx context.Context, teamID int32, nationality, image, bio string) error {
	err := db.Sql.UpdateTeam(ctx, db.UpdateTeamParams{
		ID:          teamID,
		Nationality: sql.NullString{String: nationality, Valid: nationality != ""},
		Image:       sql.NullString{String: image, Valid: image != ""},
		Bio:         sql.NullString{String: bio, Valid: bio != ""},
	})
	if err != nil {
		return err
	}

	return nil
}
