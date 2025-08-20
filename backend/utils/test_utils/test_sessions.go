package test_utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"trxd/utils"

	"github.com/gofiber/fiber/v2"
)

type apiTestSession struct {
	t        *testing.T
	app      *fiber.App
	Cookies  []*http.Cookie
	lastResp *http.Response
	UserID   *int32
}

func NewApiTestSession(t *testing.T, app *fiber.App) *apiTestSession {
	return &apiTestSession{
		t:       t,
		app:     app,
		Cookies: []*http.Cookie{},
	}
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

func (s *apiTestSession) Request(method string, url string, body interface{}, expectedStatus int) *http.Response {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			s.t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	if url[0] != '/' {
		url = "/" + url
	}

	req, err := http.NewRequest(method, "/api"+url, bytes.NewReader(reqBody))
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

	if expectedStatus != -1 && resp.StatusCode != expectedStatus {
		s.t.Errorf("%s %s: Expected status %d, got %d", method, url, expectedStatus, resp.StatusCode)
	}

	s.updateCookies(resp.Cookies())

	s.lastResp = resp
	return resp
}

func (s *apiTestSession) Get(url string, body interface{}, expectedStatus int) *http.Response {
	return s.Request(http.MethodGet, url, body, expectedStatus)
}

func (s *apiTestSession) Post(url string, body interface{}, expectedStatus int) *http.Response {
	return s.Request(http.MethodPost, url, body, expectedStatus)
}

func (s *apiTestSession) Put(url string, body interface{}, expectedStatus int) *http.Response {
	return s.Request(http.MethodPut, url, body, expectedStatus)
}

func (s *apiTestSession) Patch(url string, body interface{}, expectedStatus int) *http.Response {
	return s.Request(http.MethodPatch, url, body, expectedStatus)
}

func (s *apiTestSession) Delete(url string, body interface{}, expectedStatus int) *http.Response {
	return s.Request(http.MethodDelete, url, body, expectedStatus)
}

func (s *apiTestSession) Body() interface{} {
	defer s.lastResp.Body.Close()
	bodyBytes, err := io.ReadAll(s.lastResp.Body)
	if err != nil {
		s.t.Fatalf("Failed to read response body: %v", err)
	}

	var jsonDecoded interface{}
	err = json.Unmarshal(bodyBytes, &jsonDecoded)
	if err != nil {
		return nil
	}

	return jsonDecoded
}

func (s *apiTestSession) CheckResponse(expectedResponse interface{}) {
	jsonDecoded := s.Body()

	err := utils.Compare(expectedResponse, jsonDecoded)
	if err != nil {
		s.t.Fatalf("Response body does not match: %v", err)
	}
}
