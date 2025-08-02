package tests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

var testRegister = []struct {
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
		testBody:         JSON{"username": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", consts.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.ShortPassword),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"username": strings.Repeat("a", consts.MaxNameLength+1), "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongName),
	},
	{
		testBody:         JSON{"username": "test", "email": strings.Repeat("a", consts.MaxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongEmail),
	},
	{
		testBody:         JSON{"username": "test", "email": "invalid-email", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidEmail),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.UserAlreadyExists),
	},
	{
		testBody:         JSON{"username": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
}

func TestRegister(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	session := utils.NewApiTestSession(t, app)
	session.Post("/api/register", JSON{"username": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.DisabledRegistration))

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}
	session = utils.NewApiTestSession(t, app)
	session.Post("/api/register", JSON{"username": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(JSON{"username": "test", "role": string(db.UserRolePlayer)})

	for _, test := range testRegister {
		session := utils.NewApiTestSession(t, app)
		session.Post("/api/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testLogin = []struct {
	testBody         interface{}
	register         bool
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"email": strings.Repeat("a", consts.MaxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongEmail),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusUnauthorized,
		expectedResponse: errorf(consts.InvalidCredentials),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:         true,
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": string(db.UserRolePlayer)},
	},
}

func TestLogin(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	for _, test := range testLogin {
		if test.register {
			session := utils.NewApiTestSession(t, app)
			session.Post("/api/register", test.testBody, http.StatusOK)
		}

		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testLogout = []struct {
	testBody       interface{}
	register       bool
	login          bool
	expectedStatus int
}{
	{
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:       true,
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"email": "test@test.test", "password": "testpass"},
		login:          true,
		expectedStatus: http.StatusOK,
	},
}

func TestLogout(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	for _, test := range testLogout {
		session := utils.NewApiTestSession(t, app)

		if test.register {
			session.Post("/api/register", test.testBody, http.StatusOK)
			session.Post("/api/player/register-team", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
		} else if test.login {
			session.Post("/api/login", test.testBody, http.StatusOK)
		}
		session.Post("/api/logout", test.testBody, test.expectedStatus)

		for _, cookie := range session.Cookies {
			if cookie.Name == "session_id" && cookie.Value != "" {
				t.Errorf("Expected session_id cookie to be cleared, got %s", cookie.Value)
			}
		}

		session.Post("/api/player/register-team", JSON{"name": "test-team2", "password": "testpass"}, http.StatusUnauthorized)
	}
}

var testUpdateUser = []struct {
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
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxNameLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongName),
	},
	{
		testBody:         JSON{"nationality": strings.Repeat("a", consts.MaxNationalityLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongNationality),
	},
	{
		testBody:       JSON{"name": "a", "nationality": "a", "image": "a"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "b", "nationality": "b"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "c"},
		expectedStatus: http.StatusOK,
	},
}

func TestUpdateUser(t *testing.T) {
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

	for _, test := range testUpdateUser {
		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/api/player/user", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testResetUserPassword = []struct {
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
		testBody:         JSON{"user_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidUserID),
	},
	{
		testBody:       JSON{"user_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"user_id": 0},
		expectedStatus: http.StatusOK,
	},
}

func TestResetUserPassword(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

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
	password := "testpass"

	for _, test := range testResetUserPassword {
		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["user_id"]; ok && content == 0 {
				test.testBody.(JSON)["user_id"] = user.ID
			}
		}
		session.Post("/api/admin/reset-user-password", test.testBody, test.expectedStatus)
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
		session.Post("/api/login", JSON{"email": "test@test.test", "password": password}, http.StatusOK)
		session.CheckResponse(JSON{"username": "test", "role": string(db.UserRolePlayer)})
	}
}
