package configs_update_test

import (
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
	test_utils.Main(m, "../../../", "configs_update")
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

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRoleAdmin)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/configs/update", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
