package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/crypto_utils"

	"github.com/lib/pq"
	"github.com/tde-nico/log"
)

var db *sql.DB
var Sql *sqlc.Queries

func init() {
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		log.Fatal("Failed to set timezone:", err)
	}
}

func ConnectDB(info *utils.DBInfo, test ...bool) error {
	var err error
	db, err = sql.Open("postgres", info.ConnectionString())
	if err != nil {
		return err
	}

	if info.MaxConnections <= 0 {
		log.Fatal("invalid max connections: must be greater than 0")
	} else if info.MaxConnections > 100 {
		log.Warn("max connections is set to a high value, hard cap set to 100")
		info.MaxConnections = 100
	}

	db.SetMaxOpenConns(info.MaxConnections)
	db.SetMaxIdleConns(info.MaxConnections)
	db.SetConnMaxIdleTime(time.Hour)

	Sql = sqlc.New(db)

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
	if Sql != nil {
		Sql.Close()
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
	"allow-register":            false,
	"chall-min-points":          50,
	"chall-points-decay":        15,
	"instance-lifetime":         30 * 60, // 30 minutes
	"reclaim-instance-interval": 5 * 60,  // 5 minutes
	"instance-max-memory":       512,
	"instance-max-cpu":          "1.0",
	"min-port":                  20000,
	"max-port":                  30000,
	"hash-len":                  12,
	"secret":                    "",
	"domain":                    "",
	"discord-webhook":           "",
}

func InitConfigs() error {
	if secret, ok := defaultConfigs["secret"]; ok && secret == "" {
		randSecret, err := crypto_utils.GeneratePassword()
		if err != nil {
			return fmt.Errorf("failed to generate random secret: %v", err)
		}
		defaultConfigs["secret"] = randSecret
	}

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

	success, err = ExecSQLFile("sql/functions.sql")
	if err != nil {
		return false, err
	}
	if !success {
		return false, nil
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

func BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
