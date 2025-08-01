package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

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

func Test404(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	session := utils.NewApiTestSession(t, app)
	session.Get("/nonexistent-endpoint", errorf(consts.EndpointNotFound), http.StatusNotFound)
}
