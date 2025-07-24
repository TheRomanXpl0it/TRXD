package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"trxd/db"

	"github.com/gofiber/fiber/v2"
)

func TestMain(m *testing.M) {
	if err := db.OpenTestDB("test_api"); err != nil {
		fmt.Printf("Failed to open test database: %v\n", err)
		os.Exit(1)
	}
	defer db.CloseTestDB()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func apiRequest(app *fiber.App, method string, url string, body interface{}, cookies []*http.Cookie) (*http.Response, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	r, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}

	return app.Test(r)
}

func checkApiResponse(resp *http.Response, expectedStatus int, expectedError string) error {
	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var jsonDecoded map[string]string
	if expectedError != "" {
		err = json.Unmarshal(bodyBytes, &jsonDecoded)
		if err != nil {
			return err
		}
		jsonError, ok := jsonDecoded["error"]
		if !ok {
			return fmt.Errorf("Expected error field in response, got: %s", bodyBytes)
		}
		if jsonError != expectedError {
			return fmt.Errorf("Expected error '%s', got '%s'", expectedError, jsonError)
		}
	}
	return nil
}
