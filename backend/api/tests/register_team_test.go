package tests

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
)

var testRegisterTeam = []struct {
	testBody         interface{}
	secondUser       bool
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.InvalidJSON),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", api.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.ShortPassword),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", api.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongPassword),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", api.MaxNameLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(api.LongName),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test"},
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(api.AlreadyInTeam),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		secondUser:       true,
		expectedResponse: errorf(api.TeamAlreadyExists),
	},
	{
		testBody:         JSON{"name": "test1", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		secondUser:       true,
		expectedResponse: JSON{"name": "test1"},
	},
}

func TestRegisterTeam(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	session := newApiTestSession(t, app)
	session.Post("/register", JSON{
		"email":    "test@test.test",
		"username": "test",
		"password": "testpass",
	}, http.StatusOK)
	session.Post("/register", JSON{
		"email":    "test2@test.test",
		"username": "test2",
		"password": "testpass",
	}, http.StatusOK)

	for _, test := range testRegisterTeam {
		session := newApiTestSession(t, app)
		if test.secondUser {
			session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		}
		session.Post("/register-team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
