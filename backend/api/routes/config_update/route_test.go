package config_update_test

import (
	"context"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "config_update")
}

var testConfigUpdate = []struct {
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
		testBody:         JSON{"key": "allow-register"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"value": "true"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"key": "aaaaaaaaa", "value": "true"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ConfigNotFound),
	},
	{
		testBody:       JSON{"key": "allow-register", "value": "false"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "allow-register", "value": "true"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "allow-register", "value": "true"},
		expectedStatus: http.StatusOK,
	},
}

func TestConfigUpdate(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	user, err := user_register.RegisterUser(context.Background(), "test", "test@test.test", "testpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	for _, test := range testConfigUpdate {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/config", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
