package api

import (
	"net/http"
	"strings"
	"testing"
	"trxd/db"
)

var testRegisterTeam = []struct {
	testBody         interface{}
	secondUser       bool
	expectedStatus   int
	expectedResponse map[string]interface{}
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": invalidJSON},
	},
	{
		testBody:         map[string]string{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"name": "test", "password": strings.Repeat("a", minPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": shortPassword},
	},
	{
		testBody:         map[string]string{"name": "test", "password": strings.Repeat("a", maxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longPassword},
	},
	{
		testBody:         map[string]string{"name": strings.Repeat("a", maxNameLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longName},
	},
	{
		testBody:         map[string]string{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: map[string]interface{}{"name": "test"},
	},
	{
		testBody:         map[string]string{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: map[string]interface{}{"error": alreadyInTeam},
	},
	{
		testBody:         map[string]string{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		secondUser:       true,
		expectedResponse: map[string]interface{}{"error": teamAlreadyExists},
	},
	{
		testBody:         map[string]string{"name": "test1", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		secondUser:       true,
		expectedResponse: map[string]interface{}{"name": "test1"},
	},
}

func TestRegisterTeam(t *testing.T) {
	db.DeleteAll()
	app := SetupApp()
	defer app.Shutdown()

	session := newApiTestSession(t, app)
	session.Request(http.MethodPost, "/register", map[string]string{
		"email":    "test@test.test",
		"username": "test",
		"password": "testpass",
	}, http.StatusOK)
	session.Request(http.MethodPost, "/register", map[string]string{
		"email":    "test2@test.test",
		"username": "test2",
		"password": "testpass",
	}, http.StatusOK)

	for _, test := range testRegisterTeam {
		session := newApiTestSession(t, app)
		if test.secondUser {
			session.Request(http.MethodPost, "/login", map[string]string{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Request(http.MethodPost, "/login", map[string]string{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		}
		session.Request(http.MethodPost, "/register-team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
