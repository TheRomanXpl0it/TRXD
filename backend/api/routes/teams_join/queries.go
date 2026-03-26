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

	valid, err := crypto_utils.Verify(password, team.PasswordSalt, team.PasswordHash)
	if err != nil || !valid {
		return nil, err
	}

	return team, nil
}

func AddTeamMember(ctx context.Context, teamID int32, userID int32) error {
	err := db.Sql.AddTeamMember(ctx, sqlc.AddTeamMemberParams{
		TeamID: sql.NullInt32{Int32: teamID, Valid: true},
		ID:     userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func JoinTeam(ctx context.Context, name, password string, userID int32) (*sqlc.Team, error) {
	team, err := authTeam(ctx, name, password)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	err = AddTeamMember(ctx, team.ID, userID)
	if err != nil {
		return nil, err
	}

	return team, nil
}
