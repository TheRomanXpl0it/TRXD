package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var test_db *sql.DB

func OpenTestDB() error {
	err := os.Chdir("..")
	if err != nil {
		return err
	}

	err = godotenv.Load(".env")
	if err != nil {
		return err
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf(connStrTemplate, user, password, dbName)
	test_db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = test_db.Exec(`DROP DATABASE IF EXISTS test_db;`)
	if err != nil {
		return err
	}
	_, err = test_db.Exec(`CREATE DATABASE test_db;`)
	if err != nil {
		return err
	}
	defer test_db.Exec(`DROP DATABASE IF EXISTS test_db;`)

	err = setupTestDB()
	if err != nil {
		return err
	}

	return nil
}

func setupTestDB() error {
	err := ConnectDB(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"test_db",
	)
	if err != nil {
		return err
	}

	err = ExecSQLFile("sql/schema.sql")
	if err != nil {
		return err
	}
	files, err := os.ReadDir("sql/triggers")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		err = ExecSQLFile("sql/triggers/" + file.Name())
		if err != nil {
			return err
		}
	}
	err = ExecSQLFile("sql/tests.sql")
	if err != nil {
		return err
	}

	return nil
}

func CloseTestDB() {
	CloseDB()
	if test_db != nil {
		test_db.Close()
	}
}
