package test_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"trxd/utils"

	"github.com/gofiber/fiber/v2"
)

type apiTestSession struct {
	t        *testing.T
	app      *fiber.App
	global   bool
	Cookies  []*http.Cookie
	lastResp *http.Response
}

func NewApiTestSession(t *testing.T, app *fiber.App, globalRequest ...bool) *apiTestSession {
	s := &apiTestSession{
		t:       t,
		app:     app,
		Cookies: []*http.Cookie{},
	}

	s.Get("/info", nil, http.StatusOK)

	if len(globalRequest) > 0 && globalRequest[0] {
		s.global = true
	}

	return s
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

func (s *apiTestSession) SendRequest(req *http.Request, expectedStatus int) *http.Response {
	for _, cookie := range s.Cookies {
		if cookie.Name == "csrf_" {
			req.Header.Set("X-Csrf-Token", cookie.Value)
		}
		req.AddCookie(cookie)
	}

	resp, err := s.app.Test(req)
	if err != nil {
		s.t.Fatalf("Failed to perform request: %v", err)
	}

	if expectedStatus != -1 && resp.StatusCode != expectedStatus {
		s.t.Errorf("%s %s: Expected status %d, got %d", req.Method, req.URL.Path, expectedStatus, resp.StatusCode)
	}

	s.updateCookies(resp.Cookies())

	s.lastResp = resp
	return resp
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
	if !s.global {
		url = "/api" + url
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		s.t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return s.SendRequest(req, expectedStatus)
}

func (s *apiTestSession) RequestMultipart(method string, url string, body map[string]interface{}, files []string, fileNames []string, expectedStatus int) *http.Response {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	for field, content := range body {
		switch content := content.(type) {
		case []string:
			for _, item := range content {
				fieldWriter, err := writer.CreateFormField(field)
				if err != nil {
					s.t.Fatalf("Failed to create form field %s: %v", field, err)
				}
				_, err = fieldWriter.Write([]byte(item))
				if err != nil {
					s.t.Fatalf("Failed to write form field %s: %v", field, err)
				}
			}
		default:
			fieldWriter, err := writer.CreateFormField(field)
			if err != nil {
				s.t.Fatalf("Failed to create form field %s: %v", field, err)
			}
			_, err = fmt.Fprintf(fieldWriter, "%v", content)
			if err != nil {
				s.t.Fatalf("Failed to write form field %s: %v", field, err)
			}
		}
	}

	for i, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			s.t.Fatalf("Failed to open file %s: %v", filePath, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				s.t.Fatalf("Failed to close file %s: %v", filePath, err)
			}
		}()

		fileName := filepath.Base(filePath)
		if i < len(fileNames) {
			fileName = fileNames[i]
		}
		fileWriter, err := writer.CreateFormFile("files", fileName)
		if err != nil {
			s.t.Fatalf("Failed to create form file for %s: %v", fileName, err)
		}
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			s.t.Fatalf("Failed to copy file content for %s: %v", fileName, err)
		}
	}

	err := writer.Close()
	if err != nil {
		s.t.Fatalf("Failed to close multipart writer: %v", err)
	}

	if url[0] != '/' {
		url = "/" + url
	}
	if !s.global {
		url = "/api" + url
	}

	req, err := http.NewRequest(method, url, &requestBody)
	if err != nil {
		s.t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return s.SendRequest(req, expectedStatus)
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

func (s *apiTestSession) Options(url string, body interface{}, expectedStatus int) *http.Response {
	return s.Request(http.MethodOptions, url, body, expectedStatus)
}

func (s *apiTestSession) GetMultipart(url string, body map[string]interface{}, files []string, expectedStatus int) *http.Response {
	return s.RequestMultipart(http.MethodGet, url, body, files, nil, expectedStatus)
}

func (s *apiTestSession) PostMultipart(url string, body map[string]interface{}, files []string, expectedStatus int) *http.Response {
	return s.RequestMultipart(http.MethodPost, url, body, files, nil, expectedStatus)
}

func (s *apiTestSession) PutMultipart(url string, body map[string]interface{}, files []string, expectedStatus int) *http.Response {
	return s.RequestMultipart(http.MethodPut, url, body, files, nil, expectedStatus)
}

func (s *apiTestSession) PatchMultipart(url string, body map[string]interface{}, files []string, expectedStatus int) *http.Response {
	return s.RequestMultipart(http.MethodPatch, url, body, files, nil, expectedStatus)
}

func (s *apiTestSession) DeleteMultipart(url string, body map[string]interface{}, files []string, expectedStatus int) *http.Response {
	return s.RequestMultipart(http.MethodDelete, url, body, files, nil, expectedStatus)
}

func (s *apiTestSession) Body() interface{} {
	defer func() {
		err := s.lastResp.Body.Close()
		if err != nil {
			s.t.Fatalf("Failed to close response body: %v", err)
		}
	}()
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
		Fatalf(s.t, "Response body does not match: %v", err)
	}
}
