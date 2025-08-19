package test_utils

import (
	"context"
	"fmt"
	"os"
	"testing"
	"trxd/api/routes/users_register"
	"trxd/db"
	"trxd/db/sqlc"
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

	err = db.OpenTestDB("test_" + name)
	if err != nil {
		fatalf("Failed to open test database: %v\n", err)
	}
	defer db.CloseTestDB()

	err = db.DeleteAll()
	if err != nil {
		fatalf("Failed to delete all data: %v\n", err)
	}

	err = db.InitConfigs()
	if err != nil {
		fatalf("Failed to initialize configs: %v\n", err)
	}

	err = db.InsertMockData()
	if err != nil {
		fatalf("Failed to insert mock data: %v\n", err)
	}

	err = db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		fatalf("Failed to update config: %v\n", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func RegisterUser(t *testing.T, name, email, password string, role sqlc.UserRole) *sqlc.User {
	user, err := users_register.RegisterUser(t.Context(), name, email, password, role)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("Registered user is nil")
	}

	return user
}
