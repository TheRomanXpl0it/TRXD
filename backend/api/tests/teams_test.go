package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

var testRegisterTeam = []struct {
	testBody         interface{}
	secondUser       bool
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", consts.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.ShortPassword),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxNameLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongName),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test"},
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.AlreadyInTeam),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		secondUser:       true,
		expectedResponse: errorf(consts.TeamAlreadyExists),
	},
	{
		testBody:         JSON{"name": "test1", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		secondUser:       true,
		expectedResponse: JSON{"name": "test1"},
	},
}

func TestRegisterTeam(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	session := utils.NewApiTestSession(t, app)
	session.Post("/api/register", JSON{
		"username": "test",
		"email":    "test@test.test",
		"password": "testpass",
	}, http.StatusOK)
	session.Post("/api/register", JSON{
		"username": "test2",
		"email":    "test2@test.test",
		"password": "testpass",
	}, http.StatusOK)

	for _, test := range testRegisterTeam {
		session := utils.NewApiTestSession(t, app)
		if test.secondUser {
			session.Post("/api/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/api/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		}
		session.Post("/api/player/register-team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testJoinTeam = []struct {
	testBody         interface{}
	secondUser       bool
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", consts.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.ShortPassword),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxNameLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongName),
	},
	{
		testBody:         JSON{"name": "test1", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.InvalidTeamCredentials),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpassa"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.InvalidTeamCredentials),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test"},
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		secondUser:       true,
		expectedResponse: errorf(consts.AlreadyInTeam),
	},
}

func TestJoinTeam(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	session := utils.NewApiTestSession(t, app)
	session.Post("/api/register", JSON{
		"username": "test",
		"email":    "test@test.test",
		"password": "testpass",
	}, http.StatusOK)
	session.Post("/api/player/register-team", JSON{
		"name":     "test",
		"password": "testpass",
	}, http.StatusOK)
	session.Post("/api/register", JSON{
		"username": "test2",
		"email":    "test2@test.test",
		"password": "testpass",
	}, http.StatusOK)

	for _, test := range testJoinTeam {
		session := utils.NewApiTestSession(t, app)
		if test.secondUser {
			session.Post("/api/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/api/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
		}
		session.Post("/api/player/join-team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testUpdateTeam = []struct {
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
		testBody:         JSON{"nationality": strings.Repeat("a", consts.MaxNationalityLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongNationality),
	},
	{
		testBody:       JSON{"nationality": "a", "image": "a", "bio": "a"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"nationality": "b", "image": "b"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"nationality": "c"},
		expectedStatus: http.StatusOK,
	},
}

func TestUpdateTeam(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	session := utils.NewApiTestSession(t, app)
	session.Post("/api/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/api/player/register-team", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)

	for _, test := range testUpdateTeam {
		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/api/player/team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
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
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	admin, err := db.RegisterUser(context.Background(), "admin", "admin@test.test", "adminpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if admin == nil {
		t.Fatal("User registration returned nil")
	}

	user, err := db.RegisterUser(context.Background(), "test", "test@test.test", "testpass")
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	team, err := db.RegisterTeam(context.Background(), "test", "testpass", user.ID)
	if err != nil {
		t.Fatalf("Failed to register test team: %v", err)
	}
	if team == nil {
		t.Fatal("Team registration returned nil")
	}
	password := "testpass"

	for i, test := range testResetTeamPassword {
		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["team_id"]; ok && content == 0 {
				test.testBody.(JSON)["team_id"] = team.ID
			}
		}
		session.Post("/api/admin/reset-team-password", test.testBody, test.expectedStatus)
		if test.expectedStatus == http.StatusOK {
			body := session.Body()
			newPasswordInterface, ok := body["new_password"]
			if !ok {
				t.Fatalf("Expected 'new_password' in response, got: %v", body)
			}
			password, ok = newPasswordInterface.(string)
			if !ok {
				t.Fatalf("Expected 'new_password' to be a string, got: %T", newPasswordInterface)
			}
		}

		session = utils.NewApiTestSession(t, app)
		session.Post("/api/register", JSON{"username": "test", "email": fmt.Sprintf("test%d@test.test", i), "password": "testpass"}, http.StatusOK)
		session.Post("/api/player/join-team", JSON{"name": "test", "password": password}, http.StatusOK)
		session.Post("/api/player/submit", JSON{}, http.StatusBadRequest)
	}
}
