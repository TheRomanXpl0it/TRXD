package utils

import (
	"fmt"
	"os"
	"strconv"
)

type DBInfo struct {
	User           string
	Password       string
	DBName         string
	Host           string
	Port           int
	MaxConnections int
}

const connStrTemplate = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

func (info *DBInfo) ConnectionString() string {
	return fmt.Sprintf(connStrTemplate, info.User, info.Password, info.Host, info.Port, info.DBName)
}

func GetDBInfoFromEnv() (*DBInfo, error) {
	portStr := os.Getenv("POSTGRES_PORT")
	if portStr == "" {
		portStr = "5432"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_PORT: %v", err)
	}
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	maxConnsStr := os.Getenv("MAX_CONNECTIONS")
	if maxConnsStr == "" {
		maxConnsStr = "50"
	}
	maxConns, err := strconv.Atoi(maxConnsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_CONNECTIONS: %v", err)
	}

	info := DBInfo{
		User:           os.Getenv("POSTGRES_USER"),
		Password:       os.Getenv("POSTGRES_PASSWORD"),
		DBName:         os.Getenv("POSTGRES_DB"),
		Host:           host,
		Port:           port,
		MaxConnections: maxConns,
	}

	if info.User == "" || info.Password == "" || info.DBName == "" {
		return nil, fmt.Errorf("POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set")
	}

	return &info, nil
}
