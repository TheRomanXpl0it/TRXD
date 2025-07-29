package tests

import (
	"fmt"
	"os"
	"testing"
	"trxd/db"
)

func TestMain(m *testing.M) {
	err := os.Chdir("../../")
	if err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		os.Exit(1)
	}

	if err := db.OpenTestDB("test_api"); err != nil {
		fmt.Printf("Failed to open test database: %v\n", err)
		os.Exit(1)
	}
	defer db.CloseTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}
