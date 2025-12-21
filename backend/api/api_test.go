package api_test

import (
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
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

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
	api.SetupApi(t.Context(), app)
	defer api.Shutdown(app)

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

// 	app := api.SetupApp(t.Context())
// 	defer api.Shutdown(app)

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
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

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
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	session := test_utils.NewApiTestSession(t, app, true)
	session.Get("/static/invalid_file", nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.NotFound))
	session.Get("/static/countries.json", nil, http.StatusOK)
}

func TestUserMode(t *testing.T) {
	test_utils.UpdateConfig(t, "user-mode", "true")
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	enpoints := []struct {
		method string
		route  string
	}{
		{method: http.MethodGet, route: "/users"},
		{method: http.MethodGet, route: "/users/0"},
		{method: http.MethodPost, route: "/teams/register"},
		{method: http.MethodPost, route: "/teams/join"},
		{method: http.MethodPatch, route: "/teams"},
		{method: http.MethodPatch, route: "/teams/password"},
	}
	session := test_utils.NewApiTestSession(t, app)
	for _, ep := range enpoints {
		session.Request(ep.method, ep.route, nil, http.StatusNotFound)
		session.CheckResponse(errorf(consts.NotFound))
	}
}
