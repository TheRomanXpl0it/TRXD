package configs_get

import (
	"context"
	"trxd/db"
)

type Config struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func GetConfigs(ctx context.Context) ([]Config, error) {
	configs, err := db.Sql.GetConfigs(ctx)
	if err != nil {
		return nil, err
	}

	resConfigs := make([]Config, len(configs))
	for i, config := range configs {
		resConfigs[i] = Config{
			Key:   config.Key,
			Type:  config.Type,
			Value: config.Value,
		}
		if config.Description.Valid {
			resConfigs[i].Description = config.Description.String
		}
	}

	return resConfigs, nil
}
