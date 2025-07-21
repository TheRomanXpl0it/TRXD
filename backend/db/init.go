package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

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
	queries.Close()
	db.Close()
}
