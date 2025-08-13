package user_logout_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "user_logout")
}

var testUserLogout = []struct {
	testBody       interface{}
	register       bool
	login          bool
	expectedStatus int
}{
	{
		expectedStatus: http.StatusUnauthorized,
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

func TestUserLogout(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testUserLogout {
		session := test_utils.NewApiTestSession(t, app)

		if test.register {
			session.Post("/register", test.testBody, http.StatusOK)
			session.Post("/teams", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
		} else if test.login {
			session.Post("/login", test.testBody, http.StatusOK)
		}
		session.Post("/logout", test.testBody, test.expectedStatus)

		for _, cookie := range session.Cookies {
			if cookie.Name == "session_id" && cookie.Value != "" {
				t.Errorf("Expected session_id cookie to be cleared, got %s", cookie.Value)
			}
		}

		session.Post("/teams", JSON{"name": "test-team2", "password": "testpass"}, http.StatusUnauthorized)
	}
}
