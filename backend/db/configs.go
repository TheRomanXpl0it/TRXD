package db

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

func CreateConfig(ctx context.Context, key string, value any) (bool, error) {
	val := fmt.Sprint(value)

	err := Sql.CreateConfig(ctx, sqlc.CreateConfigParams{
		Key:   key,
		Type:  fmt.Sprintf("%T", value),
		Value: val,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return false, nil
			}
		}
		return false, err
	}

	if RedisStorage == nil {
		return true, nil
	}

	if val == "" {
		val = "\x00"
	}
	err = RedisStorage.SetWithContext(ctx, key, []byte(val), 0)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateConfig(ctx context.Context, key string, value any) error {
	val := fmt.Sprint(value)

	err := Sql.UpdateConfig(ctx, sqlc.UpdateConfigParams{
		Key:   key,
		Value: val,
	})
	if err != nil {
		return err
	}

	if RedisStorage == nil {
		return nil
	}

	if val == "" {
		val = "\x00"
	}
	err = RedisStorage.SetWithContext(ctx, key, []byte(val), 0)
	if err != nil {
		return err
	}

	return nil
}

func GetCompleteConfig(ctx context.Context, key string) (*sqlc.Config, error) {
	config, err := Sql.GetConfig(ctx, key)
	if err != nil {
		if err = sql.ErrNoRows; err != nil {
			return nil, nil
		}
		return nil, err
	}

	if RedisStorage == nil {
		return &config, nil
	}

	val := config.Value
	if val == "" {
		val = "\x00"
	}
	err = RedisStorage.SetWithContext(ctx, key, []byte(val), 0)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfig(ctx context.Context, key string) (string, error) {
	if RedisStorage != nil {
		val, err := RedisStorage.GetWithContext(ctx, key)
		if err != nil {
			return "", err
		}
		if val != nil {
			if string(val) == "\x00" {
				return "", nil
			}
			return string(val), nil
		}
	}

	config, err := GetCompleteConfig(ctx, key)
	if err != nil {
		return "", err
	}
	if config == nil {
		return "", nil
	}

	return config.Value, nil
}
