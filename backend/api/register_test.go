package api

import (
	"bytes"
	"encoding/json"
	"io"
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
		body := test.testBody
		var reqBody []byte
		if body != nil {
			var err error
			reqBody, err = json.Marshal(body)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
		}

		r, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		r.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(r)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != test.expectedStatus {
			t.Errorf("Expected status %d, got %d", test.expectedStatus, resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		var jsonDecoded map[string]string
		if test.expectedError != "" {
			err = json.Unmarshal(bodyBytes, &jsonDecoded)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}
			jsonError, ok := jsonDecoded["error"]
			if !ok {
				t.Fatalf("Expected error field in response, got: %s", bodyBytes)
			}
			if jsonError != test.expectedError {
				t.Errorf("Expected error '%s', got '%s'", test.expectedError, jsonError)
			}
		}
	}
}
