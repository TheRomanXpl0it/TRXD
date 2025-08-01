package tests

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
)

var testLogout = []struct {
	testBody       interface{}
	register       bool
	login          bool
	expectedStatus int
}{
	{
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:       true,
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"email": "test@test.test", "password": "testpass"},
		login:          true,
		expectedStatus: http.StatusOK,
	},
}

//TODO: test auth required endpoints

func TestLogout(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testLogout {
		session := utils.NewApiTestSession(t, app)

		if test.register {
			session.Post("/register", test.testBody, http.StatusOK)
			session.Post("/register-team", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
		} else if test.login {
			session.Post("/login", test.testBody, http.StatusOK)
		}
		session.Post("/logout", test.testBody, test.expectedStatus)

		for _, cookie := range session.Cookies {
			if cookie.Name == "session_id" && cookie.Value != "" {
				t.Errorf("Expected session_id cookie to be cleared, got %s", cookie.Value)
			}
		}

		session.Post("/register-team", JSON{"name": "test-team2", "password": "testpass"}, http.StatusUnauthorized)
	}
}
