package challenges_update

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/db"
	"trxd/db/sqlc"
)

func nullString(src *string) sql.NullString {
	if src == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *src, Valid: true}
}

func nullBool(src *bool) sql.NullBool {
	if src == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *src, Valid: true}
}

func nullInt32(src *int32) sql.NullInt32 {
	if src == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *src, Valid: true}
}

func nullDeployType(src *sqlc.DeployType) sqlc.NullDeployType {
	if src == nil {
		return sqlc.NullDeployType{Valid: false}
	}
	return sqlc.NullDeployType{DeployType: *src, Valid: true}
}

func nullScoreType(src *sqlc.ScoreType) sqlc.NullScoreType {
	if src == nil {
		return sqlc.NullScoreType{Valid: false}
	}
	return sqlc.NullScoreType{ScoreType: *src, Valid: true}
}

func IsChallEmpty(data *UpdateChallParams) bool {
	if data.Name == "" && data.Category == "" && data.Description == nil && data.Difficulty == nil &&
		data.Authors == nil && data.Tags == nil && data.Type == nil && data.Hidden == nil && data.MaxPoints == nil &&
		data.ScoreType == nil && data.Host == nil && data.Port == nil {
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
		ChallID:     *data.ChallID,
		Name:        sql.NullString{String: data.Name, Valid: data.Name != ""},
		Category:    sql.NullString{String: data.Category, Valid: data.Category != ""},
		Description: nullString(data.Description),
		Difficulty:  nullString(data.Difficulty),
		Type:        nullDeployType(data.Type),
		Hidden:      nullBool(data.Hidden),
		MaxPoints:   nullInt32(data.MaxPoints),
		ScoreType:   nullScoreType(data.ScoreType),
		Host:        nullString(data.Host),
		Port:        nullInt32(data.Port),
	}

	if data.Authors != nil {
		challParams.Authors = *data.Authors
	}
	if data.Tags != nil {
		challParams.Tags = *data.Tags
	}

	dockerParams := sqlc.UpdateDockerConfigsParams{
		ChallID:    *data.ChallID,
		Image:      nullString(data.Image),
		Compose:    nullString(data.Compose),
		HashDomain: nullBool(data.HashDomain),
		Lifetime:   nullInt32(data.Lifetime),
		Envs:       nullString(data.Envs),
		MaxMemory:  nullInt32(data.MaxMemory),
		MaxCpu:     nullString(data.MaxCpu),
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer db.Rollback(tx)
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
