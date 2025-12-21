package db

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := os.Chdir("../")
	if err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		os.Exit(1)
	}

	if err := OpenTestDB("test_db"); err != nil {
		fmt.Printf("Failed to open test database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		err := CloseTestDB()
		if err != nil {
			fmt.Printf("Failed to close test database: %v\n", err)
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTriggers(t *testing.T) {
	_, err := db.Exec(`SELECT tests();`)
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}
}
