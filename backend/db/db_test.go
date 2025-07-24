package db

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := OpenTestDB(); err != nil {
		os.Exit(1)
	}
	defer CloseTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTriggers(t *testing.T) {
	_, err := db.Exec(`SELECT tests();`)
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}
}
