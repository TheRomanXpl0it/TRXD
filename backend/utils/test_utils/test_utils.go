package test_utils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"trxd/db"
)

func fatalf(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func Main(m *testing.M, path string, name string) {
	err := os.Chdir(path)
	if err != nil {
		fatalf("Failed to change directory: %v\n", err)
	}

	if err := db.OpenTestDB("test_api_" + name); err != nil {
		fatalf("Failed to open test database: %v\n", err)
	}
	defer db.CloseTestDB()

	db.DeleteAll()
	db.InitConfigs()
	db.InsertMockData()
	err = db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		fatalf("Failed to update config: %v\n", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
