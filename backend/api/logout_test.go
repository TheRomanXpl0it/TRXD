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
	expectedError  string
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
		body := test.testBody

		var cookies []*http.Cookie
		if test.register {
			resp, err := apiRequest(app, http.MethodPost, "/register", body, nil)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d after registration, got %d", http.StatusOK, resp.StatusCode)
			}
			cookies = resp.Cookies()
		} else if test.login {
			resp, err := apiRequest(app, http.MethodPost, "/login", body, nil)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d after login, got %d", http.StatusOK, resp.StatusCode)
			}
			cookies = resp.Cookies()
		}

		resp, err := apiRequest(app, http.MethodPost, "/logout", body, cookies)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "session_id" && cookie.Value != "" {
				t.Errorf("Expected session_id cookie to be cleared, got %s", cookie.Value)
			}
		}

		err = checkApiResponse(resp, test.expectedStatus, test.expectedError)
		if err != nil {
			t.Errorf("Test failed for response: %v", err)
		}
	}
}
