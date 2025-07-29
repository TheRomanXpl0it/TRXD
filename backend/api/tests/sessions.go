package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type apiTestSession struct {
	t        *testing.T
	app      *fiber.App
	Cookies  []*http.Cookie
	lastResp *http.Response
}

func newApiTestSession(t *testing.T, app *fiber.App) *apiTestSession {
	return &apiTestSession{
		t:       t,
		app:     app,
		Cookies: []*http.Cookie{},
	}
}

func (s *apiTestSession) Request(method string, url string, body interface{}, expectedStatus int) *http.Response {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			s.t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		s.t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, cookie := range s.Cookies {
		req.AddCookie(cookie)
	}

	resp, err := s.app.Test(req)
	if err != nil {
		s.t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		s.t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	s.updateCookies(resp.Cookies())

	s.lastResp = resp
	return resp
}

func (s *apiTestSession) updateCookies(newCookies []*http.Cookie) {
	cookieMap := map[string]*http.Cookie{}
	for _, c := range s.Cookies {
		cookieMap[c.Name] = c
	}
	for _, c := range newCookies {
		cookieMap[c.Name] = c
	}

	s.Cookies = make([]*http.Cookie, 0, len(cookieMap))
	for _, c := range cookieMap {
		s.Cookies = append(s.Cookies, c)
	}
}

func (s *apiTestSession) CheckResponse(expectedResponse map[string]interface{}) {
	defer s.lastResp.Body.Close()
	bodyBytes, err := io.ReadAll(s.lastResp.Body)
	if err != nil {
		s.t.Fatalf("Failed to read response body: %v", err)
	}

	if expectedResponse == nil {
		return
	}

	var jsonDecoded map[string]interface{}
	err = json.Unmarshal(bodyBytes, &jsonDecoded)
	if err != nil {
		s.t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(expectedResponse, jsonDecoded) {
		expectedBytes, _ := json.MarshalIndent(expectedResponse, "", "  ")
		actualBytes, _ := json.MarshalIndent(jsonDecoded, "", "  ")
		s.t.Fatalf("Response body does not match.\nExpected:\n%s\nGot:\n%s", expectedBytes, actualBytes)
	}
}
