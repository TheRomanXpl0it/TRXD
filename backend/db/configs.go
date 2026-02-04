package db

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/db/sqlc"
	"trxd/utils/consts"

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
			if pqErr.Code == consts.PGUniqueViolation {
				return false, nil
			}
		}
		return false, err
	}

	err = StorageSet(ctx, key, val)
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

	err = StorageSet(ctx, key, val)
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

	err = StorageSet(ctx, key, config.Value)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfig(ctx context.Context, key string) (string, error) {
	val, err := StorageGet(ctx, key)
	if err != nil {
		return "", err
	}
	if val != nil {
		return *val, nil
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
