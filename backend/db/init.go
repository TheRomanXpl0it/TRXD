package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
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

func ExecSQLFile(path string) error {
	if db == nil {
		return fmt.Errorf("database connection is not established")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return err
	}

	return nil
}
