package tests

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
)

var testRegister = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.InvalidJSON),
	},
	{
		testBody:         JSON{"username": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
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
		testBody:         JSON{"username": "test", "email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", api.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.ShortPassword),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", api.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongPassword),
	},
	{
		testBody:         JSON{"username": strings.Repeat("a", api.MaxNameLength+1), "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongName),
	},
	{
		testBody:         JSON{"username": "test", "email": strings.Repeat("a", api.MaxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongEmail),
	},
	{
		testBody:         JSON{"username": "test", "email": "invalid-email", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.InvalidEmail),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(api.UserAlreadyExists),
	},
	{
		testBody:         JSON{"username": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
}

func TestRegister(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testRegister {
		session := newApiTestSession(t, app)
		session.Post("/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
