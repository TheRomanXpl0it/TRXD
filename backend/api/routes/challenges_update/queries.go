package challenges_update

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
)

func IsChallEmpty(data *UpdateChallParams) bool {
	if data.Name == "" && data.Category == "" && data.Description == "" && data.Difficulty == "" &&
		len(data.Authors) == 0 && data.Type == nil && data.Hidden == nil && data.MaxPoints == nil &&
		data.ScoreType == nil && data.Host == "" && data.Port == nil && len(data.Attachments) == 0 {
		return true
	}
	return false
}

func IsDockerConfigsEmpty(data *UpdateChallParams) bool {
	if data.Image == "" && data.Compose == "" && data.HashDomain == nil && data.Lifetime == nil &&
		data.Envs == "" && data.MaxMemory == nil && data.MaxCpu == "" {
		return true
	}
	return false
}

func UpdateChallenge(ctx context.Context, data *UpdateChallParams) error {
	if data.ChallID == nil {
		return fmt.Errorf("missing challenge ID")
	}

	challParams := sqlc.UpdateChallengeParams{
		ChallID:     *data.ChallID,
		Name:        sql.NullString{String: data.Name, Valid: data.Name != ""},
		Category:    sql.NullString{String: data.Category, Valid: data.Category != ""},
		Description: sql.NullString{String: data.Description, Valid: data.Description != ""},
		Difficulty:  sql.NullString{String: data.Difficulty, Valid: data.Difficulty != ""},
		Authors:     sql.NullString{String: strings.Join(data.Authors, ","), Valid: data.Authors != nil}, // TODO: change separator
		Host:        sql.NullString{String: data.Host, Valid: data.Host != ""},
		Attachments: sql.NullString{String: strings.Join(data.Attachments, ","), Valid: data.Attachments != nil}, // TODO: change separator
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
	if data.Port != nil {
		challParams.Port = sql.NullInt32{Int32: int32(*data.Port), Valid: true}
	}

	dockerParams := sqlc.UpdateDockerConfigsParams{
		ChallID: *data.ChallID,
		Image:   sql.NullString{String: data.Image, Valid: data.Image != ""},
		Compose: sql.NullString{String: data.Compose, Valid: data.Compose != ""},
		Envs:    sql.NullString{String: data.Envs, Valid: data.Envs != ""},
		MaxCpu:  sql.NullString{String: data.MaxCpu, Valid: data.MaxCpu != ""},
	}

	if data.HashDomain != nil {
		dockerParams.HashDomain = sql.NullBool{Bool: *data.HashDomain, Valid: true}
	}
	if data.Lifetime != nil {
		dockerParams.Lifetime = sql.NullInt32{Int32: int32(*data.Lifetime), Valid: true}
	}
	if data.MaxMemory != nil {
		dockerParams.MaxMemory = sql.NullInt32{Int32: int32(*data.MaxMemory), Valid: true}
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
			// TODO: handle errors:
			//		- challenge does not exist
			//		- category does not exist
			//		- name collision
			return err
		}
	}

	if !IsDockerConfigsEmpty(data) {
		err = queries.UpdateDockerConfigs(ctx, dockerParams)
		if err != nil {
			// TODO: handle errors:
			//		- challenge does not exist
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
