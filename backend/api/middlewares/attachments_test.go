package middlewares_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

// TODO: add hidden attachment test + move this into attachments_test.go
func TestAttachments(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	module := test_utils.GetModuleName(t)
	dir := "/tmp/" + module + "/"
	test_utils.CreateDir(t, dir)
	test_utils.CreateFile(t, dir+"b.txt", "test-line 1\n")

	session := test_utils.NewApiTestSession(t, app, true)
	session.Get("/attachments/invalid_file", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	session.Post("/api/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Get("/api/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}

	var challID1, challID5 int32
	for _, chall := range body.([]interface{}) {
		switch chall.(map[string]interface{})["name"] {
		case "chall-1":
			challID1 = int32(chall.(map[string]interface{})["id"].(float64))
		case "chall-5":
			challID5 = int32(chall.(map[string]interface{})["id"].(float64))
		}
	}
	session.PatchMultipart("/api/challenges", JSON{"chall_id": challID1}, []string{dir + "b.txt"}, http.StatusOK)
	session.PatchMultipart("/api/challenges", JSON{"chall_id": challID5}, []string{dir + "b.txt"}, http.StatusOK)

	invalids := []string{
		"/invalid_file",
		"/-1",
		"/-1/",
		"/-1/a.txt",
		"/-1/a.txt/",
		"/-1/a.txt/a.txt",
		"/999",
		"/999/",
		"/999/a.txt",
		"/999/a.txt/",
		"/999/a.txt/a.txt",
		"/1",
		"/1/",
		"/1/a.txt",
		"/1/a.txt/",
		"/1/a.txt/a.txt",
		fmt.Sprintf("/%d", challID1),
		fmt.Sprintf("/%d/", challID1),
		fmt.Sprintf("/%d/a.txt", challID1),
		fmt.Sprintf("/%d/a.txt/", challID1),
		fmt.Sprintf("/%d/a.txt/a.txt", challID1),
		fmt.Sprintf("/%d/b.txt/", challID1),
		fmt.Sprintf("/%d/b.txt/b.txt", challID1),
		fmt.Sprintf("/%d", challID5),
		fmt.Sprintf("/%d/", challID5),
		fmt.Sprintf("/%d/a.txt", challID5),
		fmt.Sprintf("/%d/a.txt/", challID5),
		fmt.Sprintf("/%d/a.txt/a.txt", challID5),
		fmt.Sprintf("/%d/b.txt/", challID5),
		fmt.Sprintf("/%d/b.txt/b.txt", challID5),
	}
	for _, url := range invalids {
		session.Get("/attachments"+url, nil, http.StatusNotFound)
		session.CheckResponse(errorf(consts.NotFound))
	}

	session.Get(fmt.Sprintf("/attachments/%d/b.txt", challID1), nil, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/attachments/%d/b.txt", challID5), nil, http.StatusOK)
	session.CheckResponse(nil)

	session = test_utils.NewApiTestSession(t, app, true)
	session.Post("/api/register", JSON{"name": "test", "email": "user@email.com", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Post("/api/teams/register", JSON{"name": "test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	for _, url := range invalids {
		session.Get("/attachments"+url, nil, http.StatusNotFound)
		session.CheckResponse(errorf(consts.NotFound))
	}

	session.Get(fmt.Sprintf("/attachments/%d/b.txt", challID1), nil, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/attachments/%d/b.txt", challID5), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.NotFound))
}
