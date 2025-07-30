package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lib/pq"
	// _ "github.com/lib/pq"
)

// TODO: change the address
const connStrTemplate = "postgres://%s:%s@localhost:5432/%s?sslmode=disable"

var db *sql.DB
var queries *Queries

func ConnectDB(user string, password string, dbName string) error {
	connStr := fmt.Sprintf(connStrTemplate, user, password, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxIdleTime(time.Hour)

	queries = New(db)
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
