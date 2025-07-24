package api

import (
	"fmt"
	"os"
	"testing"
	"trxd/db"
)

func TestMain(m *testing.M) {
	if err := db.OpenTestDB("test_api"); err != nil {
		fmt.Printf("Failed to open test database: %v\n", err)
		os.Exit(1)
	}
	defer db.CloseTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}
