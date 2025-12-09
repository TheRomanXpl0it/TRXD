package challenges_update

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"
)

func IsChallEmpty(data *UpdateChallParams) bool {
	if data.Name == "" && data.Category == "" && data.Description == nil && data.Difficulty == nil &&
		data.Authors == nil && data.Type == nil && data.Hidden == nil && data.MaxPoints == nil &&
		data.ScoreType == nil && data.Host == nil && data.Port == nil && data.Attachments == nil {
		return true
	}
	return false
}

func IsDockerConfigsEmpty(data *UpdateChallParams) bool {
	if data.Image == nil && data.Compose == nil && data.HashDomain == nil && data.Lifetime == nil &&
		data.Envs == nil && data.MaxMemory == nil && data.MaxCpu == nil {
		return true
	}
	return false
}

func UpdateChallenge(ctx context.Context, data *UpdateChallParams) error {
	if data.ChallID == nil {
		return fmt.Errorf("missing challenge ID")
	}

	challParams := sqlc.UpdateChallengeParams{
		ChallID:  *data.ChallID,
		Name:     sql.NullString{String: data.Name, Valid: data.Name != ""},
		Category: sql.NullString{String: data.Category, Valid: data.Category != ""},
	}

	if data.Description != nil {
		challParams.Description = sql.NullString{String: *data.Description, Valid: true}
	}
	if data.Difficulty != nil {
		challParams.Difficulty = sql.NullString{String: *data.Difficulty, Valid: true}
	}
	if data.Authors != nil {
		challParams.Authors = *data.Authors
	}
	if data.Type != nil {
		challParams.Type = sqlc.NullDeployType{DeployType: *data.Type, Valid: true}
	}
	if data.Hidden != nil {
		challParams.Hidden = sql.NullBool{Bool: *data.Hidden, Valid: true}
	}
	if data.MaxPoints != nil {
		challParams.MaxPoints = sql.NullInt32{Int32: int32(*data.MaxPoints), Valid: true}
	}
	if data.ScoreType != nil {
		challParams.ScoreType = sqlc.NullScoreType{ScoreType: *data.ScoreType, Valid: true}
	}
	if data.Host != nil {
		challParams.Host = sql.NullString{String: *data.Host, Valid: true}
	}
	if data.Port != nil {
		challParams.Port = sql.NullInt32{Int32: int32(*data.Port), Valid: true}
	}
	if data.Attachments != nil {
		challParams.Attachments = sql.NullString{String: strings.Join(*data.Attachments, consts.Separator), Valid: true}
	}

	dockerParams := sqlc.UpdateDockerConfigsParams{
		ChallID: *data.ChallID,
	}

	if data.Image != nil {
		dockerParams.Image = sql.NullString{String: *data.Image, Valid: true}
	}
	if data.Compose != nil {
		dockerParams.Compose = sql.NullString{String: *data.Compose, Valid: true}
	}
	if data.HashDomain != nil {
		dockerParams.HashDomain = sql.NullBool{Bool: *data.HashDomain, Valid: true}
	}
	if data.Lifetime != nil {
		dockerParams.Lifetime = sql.NullInt32{Int32: int32(*data.Lifetime), Valid: true}
	}
	if data.Envs != nil {
		dockerParams.Envs = sql.NullString{String: *data.Envs, Valid: true}
	}
	if data.MaxMemory != nil {
		dockerParams.MaxMemory = sql.NullInt32{Int32: int32(*data.MaxMemory), Valid: true}
	}
	if data.MaxCpu != nil {
		dockerParams.MaxCpu = sql.NullString{String: *data.MaxCpu, Valid: true}
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	queries := db.Sql.WithTx(tx)

	if !IsChallEmpty(data) {
		err = queries.UpdateChallenge(ctx, challParams)
		if err != nil {
			return err
		}
	}

	if !IsDockerConfigsEmpty(data) {
		err = queries.UpdateDockerConfigs(ctx, dockerParams)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
