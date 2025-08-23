package users_login_test

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
		testBody:       JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:       true,
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testData {
		if test.register {
			session := test_utils.NewApiTestSession(t, app)
			session.Post("/register", test.testBody, http.StatusOK)
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", testData[len(testData)-1].testBody, http.StatusOK)
	session.CheckResponse(testData[len(testData)-1].expectedResponse)
	session.Post("/login", testData[len(testData)-1].testBody, http.StatusForbidden)
	session.CheckResponse(errorf(consts.AlreadyLoggedIn))
}
