package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"trxd/db"

	"github.com/gofiber/fiber/v2"
)

var testLogin = []struct {
	testBody       interface{}
	register       bool
	expectedStatus int
	expectedError  string
}{
	{
		testBody:       nil,
		expectedStatus: http.StatusBadRequest,
		expectedError:  invalidJSON,
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
		testBody:       map[string]string{"email": "test@test.test", "password": strings.Repeat("a", maxPasswordLength+1)},
		expectedStatus: http.StatusBadRequest,
		expectedError:  longPassword,
	},
	{
		testBody:       map[string]string{"email": strings.Repeat("a", maxEmailLength+1), "password": "testpass"},
		expectedStatus: http.StatusBadRequest,
		expectedError:  longEmail,
	},
	{
		testBody:       map[string]string{"email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusUnauthorized,
		expectedError:  invalidCredentials,
	},
	{
		testBody:       map[string]string{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:       true,
		expectedStatus: http.StatusOK,
	},
}

func loginRequest(app *fiber.App, body interface{}) (*http.Response, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	r, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	return app.Test(r)
}

func TestLogin(t *testing.T) {
	db.DeleteAll()
	app := SetupApp()
	defer app.Shutdown()

	for _, test := range testLogin {
		body := test.testBody

		if test.register {
			resp, err := apiRequest(app, http.MethodPost, "/register", body, nil)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d after registration, got %d", http.StatusOK, resp.StatusCode)
			}
		}

		resp, err := loginRequest(app, body)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		err = checkApiResponse(resp, test.expectedStatus, test.expectedError)
		if err != nil {
			t.Errorf("Test failed for response: %v", err)
		}
	}
}
