package users_register_test

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
		testBody:       JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.UserAlreadyExists),
	},
	{
		testBody:       JSON{"username": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.UpdateConfig(t, "allow-register", "false")
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.DisabledRegistration))

	test_utils.UpdateConfig(t, "allow-register", "true")
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test", "email": "allow@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Post("/users/register", JSON{"username": "test", "email": "allow+1@test.test", "password": "testpass"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.AlreadyRegistered))

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
