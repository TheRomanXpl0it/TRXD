package api

import (
	"net/http"
	"testing"
	"trxd/db"
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
		testBody:       map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:       true,
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       map[string]string{"email": "test@test.test", "password": "testpass"},
		login:          true,
		expectedStatus: http.StatusOK,
	},
}

func TestLogout(t *testing.T) {
	db.DeleteAll()
	app := SetupApp()
	defer app.Shutdown()

	for _, test := range testLogout {
		session := newApiTestSession(t, app)

		if test.register {
			session.Request(http.MethodPost, "/register", test.testBody, http.StatusOK)
		} else if test.login {
			session.Request(http.MethodPost, "/login", test.testBody, http.StatusOK)
		}
		session.Request(http.MethodPost, "/logout", test.testBody, test.expectedStatus)

		for _, cookie := range session.Cookies {
			if cookie.Name == "session_id" && cookie.Value != "" {
				t.Errorf("Expected session_id cookie to be cleared, got %s", cookie.Value)
			}
		}
	}
}
