package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/test_utils"

	"github.com/gofiber/fiber/v2"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func Test404(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/nonexistent-endpoint", nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.NotFound))
}

func TestPanic(t *testing.T) {
	consts.Testing = false
	defer func() { consts.Testing = true }()

	app := fiber.New(fiber.Config{
		AppName: consts.Name,
	})
	api.SetupFeatures(app)
	api.SetupApi(app)
	defer app.Shutdown()

	app.Get("/api/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	session := test_utils.NewApiTestSession(t, app, true)
	session.Get("/api/panic", nil, http.StatusInternalServerError)
	session.CheckResponse(errorf(consts.InternalServerError))
}

// func TestLimit(t *testing.T) {
// 	consts.Testing = false
// 	defer func() { consts.Testing = true }()

// 	app := api.SetupApp()
// 	defer app.Shutdown()

// 	session := test_utils.NewApiTestSession(t, app)
// 	for range limiter.ConfigDefault.Max - 1 {
// 		session.Get("/info", nil, http.StatusOK)
// 		session.CheckResponse(nil)
// 	}
// 	session.Get("/info", nil, http.StatusTooManyRequests)
// 	session.CheckResponse(nil)
// 	session.Get("/info", nil, http.StatusTooManyRequests)
// 	session.CheckResponse(nil)
// }

func TestCSRF(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/logout", nil, http.StatusOK)
	session.CheckResponse(nil)
	session.Post("/login", nil, http.StatusForbidden)
	session.CheckResponse(nil)
	session.Get("/info", nil, http.StatusOK)
	session.CheckResponse(nil)
	session.Post("/login", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJSON))
}

func TestStatic(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app, true)
	session.Get("/static/invalid_file", nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.NotFound))
	session.Get("/static/countries.json", nil, http.StatusOK)
}

func TestAttachments(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	module := test_utils.GetModuleName(t)
	dir := "/tmp/" + module + "/"
	test_utils.CreateDir(t, dir)
	test_utils.CreateFile(t, dir+"b.txt", "test-line 1\n")

	session := test_utils.NewApiTestSession(t, app, true)
	session.Get("/attachments/invalid_file", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.Unauthorized))

	session.Post("/api/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
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
		fmt.Sprintf("/%d", challID5),
		fmt.Sprintf("/%d/", challID5),
		fmt.Sprintf("/%d/a.txt", challID5),
		fmt.Sprintf("/%d/a.txt/", challID5),
		fmt.Sprintf("/%d/a.txt/a.txt", challID5),
		fmt.Sprintf("/%d/b.txt", challID5),
		fmt.Sprintf("/%d/b.txt/", challID5),
		fmt.Sprintf("/%d/b.txt/b.txt", challID5),
		fmt.Sprintf("/%d", challID1),
		fmt.Sprintf("/%d/", challID1),
		fmt.Sprintf("/%d/a.txt", challID1),
		fmt.Sprintf("/%d/a.txt/", challID1),
		fmt.Sprintf("/%d/a.txt/a.txt", challID1),
		fmt.Sprintf("/%d/b.txt/", challID1),
		fmt.Sprintf("/%d/b.txt/b.txt", challID1),
	}
	for _, url := range invalids {
		session.Get("/attachments"+url, nil, http.StatusNotFound)
		session.CheckResponse(errorf(consts.NotFound))
	}

	session.Get(fmt.Sprintf("/attachments/%d/b.txt", challID1), nil, http.StatusOK)
	session.CheckResponse(nil)
}
