package utils

import (
	"fmt"
	"os"
	"strconv"
)

type DBInfo struct {
	PgUser           string
	PgPassword       string
	PgDBName         string
	PgHost           string
	PgPort           int
	PgMaxConnections int
	RedisHost        string
	RedisPort        int
	RedisPassword    string
}

const connStrTemplate = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

func (info *DBInfo) ConnectionString() string {
	return fmt.Sprintf(connStrTemplate, info.PgUser, info.PgPassword, info.PgHost, info.PgPort, info.PgDBName)
}

func GetDBInfoFromEnv() (*DBInfo, error) {
	var err error
	var pgPort int
	var pgMaxConns int
	var redisPort int

	pgPortStr := os.Getenv("POSTGRES_PORT")
	if pgPortStr != "" {
		pgPort, err = strconv.Atoi(pgPortStr)
		if err != nil {
			return nil, fmt.Errorf("invalid POSTGRES_PORT: %v", err)
		}
	} else {
		pgPort = 5432
	}

	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		pgHost = "localhost"
	}

	pgMaxConnsStr := os.Getenv("MAX_CONNECTIONS")
	if pgMaxConnsStr != "" {
		pgMaxConns, err = strconv.Atoi(pgMaxConnsStr)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
		}
	} else {
		pgMaxConns = 50
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPortStr := os.Getenv("REDIS_PORT")
	if redisPortStr == "" {
		redisPort = 6379
	} else {
		redisPort, err = strconv.Atoi(redisPortStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_PORT: %v", err)
		}
	}

	info := DBInfo{
		PgUser:           os.Getenv("POSTGRES_USER"),
		PgPassword:       os.Getenv("POSTGRES_PASSWORD"),
		PgDBName:         os.Getenv("POSTGRES_DB"),
		PgHost:           pgHost,
		PgPort:           pgPort,
		PgMaxConnections: pgMaxConns,
		RedisHost:        redisHost,
		RedisPort:        redisPort,
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
	}

	if info.PgUser == "" || info.PgPassword == "" || info.PgDBName == "" {
		return nil, fmt.Errorf("POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set")
	}

	return &info, nil
}
