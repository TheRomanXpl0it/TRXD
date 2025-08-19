package users_password_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/users_register"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "user_password")
}

var testUserPassword = []struct {
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

func TestUserPassword(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	admin, err := users_register.RegisterUser(t.Context(), "admin", "admin@test.test", "adminpass", sqlc.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}
	if admin == nil {
		t.Fatal("User registration returned nil")
	}

	user, err := users_register.RegisterUser(t.Context(), "test", "test@test.test", "testpass")
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	password := "testpass"

	for _, test := range testUserPassword {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["user_id"]; ok && content == 0 {
				test.testBody.(JSON)["user_id"] = user.ID
			}
		}
		session.Patch("/users/password", test.testBody, test.expectedStatus)
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
		session.Post("/users/login", JSON{"email": "test@test.test", "password": password}, http.StatusOK)
		session.CheckResponse(nil)
	}
}
