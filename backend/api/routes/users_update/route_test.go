package users_update_test

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "user_update")
}

var testUserUpdate = []struct {
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
		testBody:         JSON{"country": strings.Repeat("a", consts.MaxCountryLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCountry),
	},
	{
		testBody:       JSON{"name": "a", "country": "a", "image": "a"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "b", "country": "b"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "c"},
		expectedStatus: http.StatusOK,
	},
}

func TestUserUpdate(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)

	for _, test := range testUserUpdate {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/users/update", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
