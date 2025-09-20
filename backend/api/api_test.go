package api_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/test_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
	session.CheckResponse(errorf(consts.EndpointNotFound))
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

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/panic", nil, http.StatusInternalServerError)
	session.CheckResponse(errorf(consts.InternalServerError))
}

func TestLimit(t *testing.T) {
	consts.Testing = false
	defer func() { consts.Testing = true }()

	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	for range limiter.ConfigDefault.Max - 1 {
		session.Get("/info", nil, http.StatusOK)
		session.CheckResponse(nil)
	}
	session.Get("/info", nil, http.StatusTooManyRequests)
	session.CheckResponse(nil)
	session.Get("/info", nil, http.StatusTooManyRequests)
	session.CheckResponse(nil)
}

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
