package challenges_update

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
)

func UpdateChallenge(ctx context.Context, data UpdateChallParams) error {
	if data.ChallID == nil {
		return fmt.Errorf("missing challenge ID")
	}

	params := sqlc.UpdateChallengeParams{
		ChallID:     *data.ChallID,
		Name:        sql.NullString{String: data.Name, Valid: data.Name != ""},
		Category:    sql.NullString{String: data.Category, Valid: data.Category != ""},
		Description: sql.NullString{String: data.Description, Valid: data.Description != ""},
		Difficulty:  sql.NullString{String: data.Difficulty, Valid: data.Difficulty != ""},
		Authors:     sql.NullString{String: strings.Join(data.Authors, ","), Valid: data.Authors != nil}, // TODO: change separator
		Host:        sql.NullString{String: data.Host, Valid: data.Host != ""},
	}

	if data.Type != nil {
		params.Type = sqlc.NullDeployType{DeployType: *data.Type, Valid: true}
	}
	if data.Hidden != nil {
		params.Hidden = sql.NullBool{Bool: *data.Hidden, Valid: true}
	}
	if data.MaxPoints != nil {
		params.MaxPoints = sql.NullInt32{Int32: int32(*data.MaxPoints), Valid: true}
	}
	if data.ScoreType != nil {
		params.ScoreType = sqlc.NullScoreType{ScoreType: *data.ScoreType, Valid: true}
	}
	if data.Port != nil {
		params.Port = sql.NullInt32{Int32: int32(*data.Port), Valid: true}
	}

	err := db.Sql.UpdateChallenge(ctx, params)
	if err != nil {
		// TODO: handle errors:
		//		- challenge does not exist
		//		- category does not exist
		return err
	}

	return nil
}

func UpdateDockerConfigs(ctx context.Context, data UpdateChallParams) error {
	if data.ChallID == nil {
		return fmt.Errorf("missing challenge ID")
	}

	params := sqlc.UpdateDockerConfigsParams{
		ChallID: *data.ChallID,
		Image:   sql.NullString{String: data.Image, Valid: data.Image != ""},
		Compose: sql.NullString{String: data.Compose, Valid: data.Compose != ""},
		Envs:    sql.NullString{String: data.Envs, Valid: data.Envs != ""},
		MaxCpu:  sql.NullString{String: data.MaxCpu, Valid: data.MaxCpu != ""},
	}

	if data.HashDomain != nil {
		params.HashDomain = sql.NullBool{Bool: *data.HashDomain, Valid: true}
	}
	if data.Lifetime != nil {
		params.Lifetime = sql.NullInt32{Int32: int32(*data.Lifetime), Valid: true}
	}
	if data.MaxMemory != nil {
		params.MaxMemory = sql.NullInt32{Int32: int32(*data.MaxMemory), Valid: true}
	}

	err := db.Sql.UpdateDockerConfigs(ctx, params)
	if err != nil {
		// TODO: handle errors:
		//		- challenge does not exist
		return err
	}

	return nil
}
