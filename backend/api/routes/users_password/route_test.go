package users_password_test

import (
	"math"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

var testData = []struct {
	isAdmin          bool
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	// Player
	{
		isAdmin:          false,
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        false,
		testBody:       JSON{},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:          false,
		testBody:         JSON{"user_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "UserID", 0)),
	},
	{
		isAdmin:          false,
		testBody:         JSON{"user_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        false,
		testBody:       JSON{"user_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        false,
		testBody:       JSON{"user_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        false,
		testBody:       JSON{"user_id": 0, "new_password": "NewPassw0rd!"},
		expectedStatus: http.StatusOK,
	},
	// Admin
	{
		isAdmin:          true,
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:          true,
		testBody:         JSON{"user_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "UserID", 0)),
	},
	{
		isAdmin:          true,
		testBody:         JSON{"user_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        true,
		testBody:       JSON{"user_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        true,
		testBody:       JSON{"user_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        true,
		testBody:       JSON{"user_id": 0, "new_password": "NewPassw0rd!"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	test_utils.RegisterUser(t, "admin", "admin@test.test", "old_adminpass", sqlc.UserRoleAdmin)
	user := test_utils.RegisterUser(t, "test", "test@test.test", "old_pass", sqlc.UserRolePlayer)
	password := "testpass"

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@test.test", "password": "old_adminpass"}, http.StatusOK)
	session.Patch("/users/password", JSON{}, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	respBody := body.(map[string]interface{})
	newAdminPassInterface, ok := respBody["new_password"]
	if !ok {
		t.Fatalf("Expected 'new_password' in response, got: %v", respBody)
	}
	AdminPass := newAdminPassInterface.(string)

	user2 := test_utils.RegisterUser(t, "test2", "test2@test.test", "testpass-2", sqlc.UserRolePlayer)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "old_pass"}, http.StatusOK)
	session.Patch("/users/password", JSON{"user_id": user2.ID, "new_password": password}, http.StatusOK)
	session.CheckResponse(nil)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": password}, http.StatusOK)
	session.CheckResponse(nil)

	for _, test := range testData {
		var email, pass string
		if test.isAdmin {
			email = "admin@test.test"
			pass = AdminPass
		} else {
			email = "test@test.test"
			pass = password
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": email, "password": pass}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["user_id"]; ok && content == 0 {
				test.testBody.(JSON)["user_id"] = user.ID
			}
		}
		session.Patch("/users/password", test.testBody, test.expectedStatus)

		var newPass interface{}
		var passOk bool
		if test.testBody != nil {
			newPass, passOk = test.testBody.(JSON)["new_password"]
		}
		if test.expectedStatus != http.StatusOK || passOk {
			session.CheckResponse(test.expectedResponse)
		}
		if test.expectedStatus == http.StatusOK {
			if passOk {
				password = newPass.(string)
			} else {
				sessionBody := session.Body()
				if sessionBody == nil {
					t.Fatal("Expected body to not be nil")
				}
				body := sessionBody.(map[string]interface{})
				newPasswordInterface, ok := body["new_password"]
				if !ok {
					t.Fatalf("Expected 'new_password' in response, got: %v", body)
				}
				password, ok = newPasswordInterface.(string)
				if !ok {
					t.Fatalf("Expected 'new_password' to be a string, got: %T", newPasswordInterface)
				}
			}
		}

		session = test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": password}, http.StatusOK)
		session.CheckResponse(nil)
	}
}
