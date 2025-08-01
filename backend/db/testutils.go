package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var test_db *sql.DB

func OpenTestDB(testDBName string) error {
	godotenv.Load(".env")

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || dbName == "" {
		return fmt.Errorf("POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set")
	}

	var err error
	connStr := fmt.Sprintf(connStrTemplate, user, password, dbName)
	test_db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = test_db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, testDBName))
	if err != nil {
		return err
	}
	_, err = test_db.Exec(fmt.Sprintf(`CREATE DATABASE %s;`, testDBName))
	if err != nil {
		return err
	}
	defer test_db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, testDBName))

	err = setupTestDB(testDBName)
	if err != nil {
		return err
	}

	return nil
}

func setupTestDB(testDBName string) error {
	err := ConnectDB(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		testDBName,
		true,
	)
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

func DeleteAll() error {
	_, err := db.Exec(`SELECT delete_all();`)
	return err
}
