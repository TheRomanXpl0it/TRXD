package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

func CreateConfig(ctx context.Context, key string, value any) (*Config, error) {
	err := queries.CreateConfig(ctx, CreateConfigParams{
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
	return &Config{
		Key:   key,
		Type:  fmt.Sprintf("%T", value),
		Value: fmt.Sprint(value),
	}, nil
}

func UpdateConfig(ctx context.Context, key string, value any) error {
	err := queries.UpdateConfig(ctx, UpdateConfigParams{
		Key:   key,
		Value: fmt.Sprint(value),
	})
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(ctx context.Context, key string) (*Config, error) {
	config, err := queries.GetConfig(ctx, key)
	if err != nil {
		if err = sql.ErrNoRows; err != nil {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}
