package team_password_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/team_register"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "team_password")
}

var testResetTeamPassword = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"team_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidUserID),
	},
	{
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
}

func TestResetTeamPassword(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(t.Context(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	admin, err := user_register.RegisterUser(t.Context(), "admin", "admin@test.test", "adminpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if admin == nil {
		t.Fatal("User registration returned nil")
	}

	user, err := user_register.RegisterUser(t.Context(), "test", "test@test.test", "testpass")
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	team, err := team_register.RegisterTeam(t.Context(), "test", "testpass", user.ID)
	if err != nil {
		t.Fatalf("Failed to register test team: %v", err)
	}
	if team == nil {
		t.Fatal("Team registration returned nil")
	}
	password := "testpass"

	for i, test := range testResetTeamPassword {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["team_id"]; ok && content == 0 {
				test.testBody.(JSON)["team_id"] = team.ID
			}
		}
		session.Patch("/teams/password", test.testBody, test.expectedStatus)
		if test.expectedStatus == http.StatusOK {
			body := session.Body().(map[string]interface{})
			newPasswordInterface, ok := body["new_password"]
			if !ok {
				t.Fatalf("Expected 'new_password' in response, got: %v", body)
			}
			password, ok = newPasswordInterface.(string)
			if !ok {
				t.Fatalf("Expected 'new_password' to be a string, got: %T", newPasswordInterface)
			}
		}

		session = test_utils.NewApiTestSession(t, app)
		session.Post("/register", JSON{"username": "test", "email": fmt.Sprintf("test%d@test.test", i), "password": "testpass"}, http.StatusOK)
		session.Put("/teams", JSON{"name": "test", "password": password}, http.StatusOK)
		session.Post("/submit", JSON{}, http.StatusBadRequest)
	}
}
