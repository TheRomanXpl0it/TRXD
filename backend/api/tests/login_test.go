package tests

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
)

var testLogin = []struct {
	testBody         interface{}
	register         bool
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.InvalidJSON),
	},
	{
		testBody:         JSON{"email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": strings.Repeat("a", api.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongPassword),
	},
	{
		testBody:         JSON{"email": strings.Repeat("a", api.MaxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongEmail),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusUnauthorized,
		expectedResponse: errorf(api.InvalidCredentials),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:         true,
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
}

func TestLogin(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testLogin {

		if test.register {
			session := newApiTestSession(t, app)
			session.Request(http.MethodPost, "/register", test.testBody, http.StatusOK)
		}

		session := newApiTestSession(t, app)
		session.Request(http.MethodPost, "/login", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
