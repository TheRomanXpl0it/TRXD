package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

// TODO: make this configurable
const connStrTemplate = "postgres://%s:%s@localhost:5432/%s?sslmode=disable"

var db *sql.DB
var queries *Queries

func ConnectDB(user string, password string, dbName string, test ...bool) error {
	connStr := fmt.Sprintf(connStrTemplate, user, password, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// TODO: make these configurable
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxIdleTime(time.Hour)

	queries = New(db)

	success, err := initDB(len(test) > 0 && test[0])
	if err != nil {
		return err
	}
	if len(test) > 0 && test[0] && !success {
		return fmt.Errorf("init already executed")
	}

	return nil
}

func CloseDB() {
	if queries != nil {
		queries.Close()
	}
	if db != nil {
		db.Close()
	}
}

func ExecSQLFile(path string) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("database connection is not established")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}

	_, err = db.Exec(string(data))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "42710" { // Object already exists error code
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

var defaultConfigs = map[string]any{
	"allow-register":     false,
	"chall-min-points":   50,
	"chall-points-decay": 15,
	"instance-lifetime":  30 * 60, // 30 minutes
	"instance-max-cpu":   "1.0",
	"instance-max-mem":   512,
}

func InitConfigs() error {
	for key, value := range defaultConfigs {
		conf, err := CreateConfig(context.Background(), key, value)
		if err != nil {
			return fmt.Errorf("failed to create config for key %s=%v: %v", key, value, err)
		}
		if conf == nil {
			return fmt.Errorf("failed to create config for key %s=%v: config already exists", key, value)
		}
	}
	return nil
}

func initDB(test ...bool) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("database connection is not established")
	}

	success, err := ExecSQLFile("sql/schema.sql")
	if err != nil {
		return false, err
	}
	if !success {
		return false, nil
	}

	err = InitConfigs()
	if err != nil {
		return false, fmt.Errorf("failed to initialize configs: %v", err)
	}

	files, err := os.ReadDir("sql/triggers")
	if err != nil {
		return false, err
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		success, err = ExecSQLFile("sql/triggers/" + file.Name())
		if err != nil {
			return false, fmt.Errorf("failed to execute trigger SQL file %s: %v", file.Name(), err)
		}
		if !success {
			return false, nil
		}
	}

	if len(test) > 0 && test[0] {
		success, err = ExecSQLFile("sql/tests.sql")
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}

	return true, nil
}
