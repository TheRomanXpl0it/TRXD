package api

import (
	"net/http"
	"strings"
	"testing"
	"trxd/db"
)

var testRegister = []struct {
	testBody       interface{}
	expectedStatus int
	expectedError  string
}{
	{
		testBody:       nil,
		expectedStatus: http.StatusBadRequest,
		expectedError:  invalidJSON,
	},
	{
		testBody:       map[string]string{"username": "test"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  missingRequiredFields,
	},
	{
		testBody:       map[string]string{"email": "test@test.test"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  missingRequiredFields,
	},
	{
		testBody:       map[string]string{"password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  missingRequiredFields,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test@test.test"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  missingRequiredFields,
	},
	{
		testBody:       map[string]string{"username": "test", "password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  missingRequiredFields,
	},
	{
		testBody:       map[string]string{"email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  missingRequiredFields,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", minPasswordLength-1)},
		expectedStatus: http.StatusBadRequest,
		expectedError:  shortPassword,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", maxPasswordLength+1)},
		expectedStatus: http.StatusBadRequest,
		expectedError:  longPassword,
	},
	{
		testBody:       map[string]string{"username": strings.Repeat("a", maxUsernameLength+1), "email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  longUsername,
	},
	{
		testBody:       map[string]string{"username": "test", "email": strings.Repeat("a", maxEmailLength+1), "password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  longEmail,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "invalid-email", "password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  invalidEmail,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusConflict,
		expectedError:  userAlreadyExists,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
}

func TestRegister(t *testing.T) {
	db.DeleteAll()
	app := SetupApp()
	defer app.Shutdown()

	for _, test := range testRegister {
		session := newApiTestSession(t, app)
		session.Request(http.MethodPost, "/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedError)
	}
}
