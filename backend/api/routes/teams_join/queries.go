package teams_join

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/crypto_utils"
)

func authTeam(ctx context.Context, name, password string) (*sqlc.Team, error) {
	team, err := db.GetTeamByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	if !crypto_utils.Verify(password, team.PasswordSalt, team.PasswordHash) {
		return nil, nil
	}

	return team, nil
}

func JoinTeam(ctx context.Context, name, password string, userID int32) (*sqlc.Team, error) {
	team, err := authTeam(ctx, name, password)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	err = db.Sql.AddTeamMember(ctx, sqlc.AddTeamMemberParams{
		TeamID: sql.NullInt32{Int32: team.ID, Valid: true},
		ID:     userID,
	})
	if err != nil {
		return nil, err
	}

	return team, nil
}
