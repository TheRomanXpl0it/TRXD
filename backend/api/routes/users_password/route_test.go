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
		expectedResponse: errorf(test_utils.Format(consts.MinError, "UserID", 0)),
	},
	{
		testBody:         JSON{"user_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
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

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	test_utils.RegisterUser(t, "admin", "admin@test.test", "adminpass", sqlc.UserRoleAdmin)
	user := test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRolePlayer)
	password := "testpass"

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["user_id"]; ok && content == 0 {
				test.testBody.(JSON)["user_id"] = user.ID
			}
		}
		session.Patch("/users/password", test.testBody, test.expectedStatus)
		if test.expectedStatus == http.StatusOK {
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
		} else {
			session.CheckResponse(test.expectedResponse)
		}

		session = test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": password}, http.StatusOK)
		session.CheckResponse(nil)
	}
}
