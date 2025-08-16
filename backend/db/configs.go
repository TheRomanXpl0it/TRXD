package db

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

func CreateConfig(ctx context.Context, key string, value any) (*sqlc.Config, error) {
	err := Sql.CreateConfig(ctx, sqlc.CreateConfigParams{
		Key:   key,
		Type:  fmt.Sprintf("%T", value),
		Value: fmt.Sprint(value),
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &sqlc.Config{
		Key:   key,
		Type:  fmt.Sprintf("%T", value),
		Value: fmt.Sprint(value),
	}, nil
}

func UpdateConfig(ctx context.Context, key string, value any) error {
	err := Sql.UpdateConfig(ctx, sqlc.UpdateConfigParams{
		Key:   key,
		Value: fmt.Sprint(value),
	})
	if err != nil {
		return err
	}

	return nil
}

func GetConfig(ctx context.Context, key string) (*sqlc.Config, error) {
	config, err := Sql.GetConfig(ctx, key)
	if err != nil {
		if err = sql.ErrNoRows; err != nil {
			return nil, nil
		}
		return nil, err
	}

	return &config, nil
}
