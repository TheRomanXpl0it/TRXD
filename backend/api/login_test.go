package api

import (
	"net/http"
	"strings"
	"testing"
	"trxd/db"
)

var testLogin = []struct {
	testBody         interface{}
	register         bool
	expectedStatus   int
	expectedResponse map[string]interface{}
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": invalidJSON},
	},
	{
		testBody:         map[string]string{"email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"email": "test@test.test", "password": strings.Repeat("a", maxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longPassword},
	},
	{
		testBody:         map[string]string{"email": strings.Repeat("a", maxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longEmail},
	},
	{
		testBody:         map[string]string{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusUnauthorized,
		expectedResponse: map[string]interface{}{"error": invalidCredentials},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:         true,
		expectedStatus:   http.StatusOK,
		expectedResponse: map[string]interface{}{"username": "test", "role": "Player"},
	},
}

func TestLogin(t *testing.T) {
	db.DeleteAll()
	app := SetupApp()
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
