package db

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestTriggers(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	err = ConnectDB(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer CloseDB()

	// TODO: add tests functions to db only from here

	_, err = db.Exec(`SELECT tests();`)
	if err != nil {
		t.Fatalf("Failed to execute tests: %v", err)
	}
}
