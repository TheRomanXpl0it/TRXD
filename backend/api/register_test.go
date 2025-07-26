package api

import (
	"net/http"
	"strings"
	"testing"
	"trxd/db"
)

var testRegister = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse map[string]interface{}
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": invalidJSON},
	},
	{
		testBody:         map[string]string{"username": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
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
		testBody:         map[string]string{"username": "test", "email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"username": "test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": missingRequiredFields},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", minPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": shortPassword},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", maxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longPassword},
	},
	{
		testBody:         map[string]string{"username": strings.Repeat("a", maxNameLength+1), "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longName},
	},
	{
		testBody:         map[string]string{"username": "test", "email": strings.Repeat("a", maxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": longEmail},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "invalid-email", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: map[string]interface{}{"error": invalidEmail},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: map[string]interface{}{"username": "test", "role": "Player"},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: map[string]interface{}{"error": userAlreadyExists},
	},
	{
		testBody:         map[string]string{"username": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: map[string]interface{}{"username": "test", "role": "Player"},
	},
}

func TestRegister(t *testing.T) {
	db.DeleteAll()
	app := SetupApp()
	defer app.Shutdown()

	for _, test := range testRegister {
		session := newApiTestSession(t, app)
		session.Request(http.MethodPost, "/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
