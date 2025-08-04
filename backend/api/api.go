package api

import (
	"encoding/json"
	"os"
	"trxd/api/auth"
	"trxd/api/routes"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/tde-nico/log"
)

func decodeJSONBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	defaultBody := string(body)

	var tmp interface{}
	if err := json.Unmarshal(body, &tmp); err == nil {
		if tmp2, err := json.MarshalIndent(tmp, "", "  "); err == nil {
			return string(tmp2)
		}
	}

	return defaultBody
}

func debugMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/countries" {
		return c.Next()
	}

	reqBody := c.BodyRaw()
	body := decodeJSONBody(reqBody)
	log.Debug("Request:", "method", c.Method(), "path", c.Path(), "body", body)

	e := c.Next()

	resStatus := c.Response().StatusCode()
	resBody := c.Response().Body()
	body = decodeJSONBody(resBody)
	log.Debug("Response:", "status", resStatus, "body", body)

	return e
}

func SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "TRXd",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:1337",
		AllowCredentials: true,
	}))

	var api fiber.Router
	if log.GetLevel() == log.DebugLevel {
		api = app.Group("/api", debugMiddleware)
	} else {
		api = app.Group("/api")
	}

	api.Get("/countries", func(c *fiber.Ctx) error { return c.JSON(consts.Countries) })

	api.Post("/register", routes.Register)
	api.Post("/login", routes.Login)
	api.Post("/logout", routes.Logout)

	api.Get("/info", auth.Spectator, routes.Info)
	api.Get("/challenges", auth.Spectator, auth.Team, routes.GetChallenges)
	api.Get("/challenges/:id", auth.Spectator, auth.Team, routes.GetChallenge)

	api.Get("/users", auth.NoAuth, routes.GetUsers)
	api.Get("/users/:id", auth.NoAuth, routes.GetUser)
	api.Patch("/users", auth.Player, routes.UpdateUser)
	api.Patch("/users/password", auth.Admin, routes.ResetUserPassword)

	api.Get("/teams", auth.NoAuth, routes.GetTeams)
	api.Get("/teams/:id", auth.NoAuth, routes.GetTeam)
	api.Post("/teams", auth.Player, routes.RegisterTeam)
	api.Put("/teams", auth.Player, routes.JoinTeam)
	api.Patch("/teams", auth.Player, auth.Team, routes.UpdateTeam)
	api.Patch("/teams/password", auth.Admin, routes.ResetTeamPassword)

	api.Post("/submit", auth.Player, auth.Team, routes.Submit)

	// api.Post("/instance", auth.Player, routes.CreateInstance)
	// api.Patch("/instance", auth.Player, routes.ExtendInstance)
	// api.Delete("/instance", auth.Player, routes.DeleteInstance)

	api.Post("/category", auth.Author, routes.CreateCategory)
	// api.Patch("/category", auth.Author, routes.UpdateCategory)
	api.Delete("/category", auth.Author, routes.DeleteCategory)

	api.Post("/challenges", auth.Author, routes.CreateChallenge)
	// api.Patch("/challenges", routes.UpdateChallenge)
	api.Delete("/challenges", auth.Author, routes.DeleteChallenge)

	api.Post("/flag", auth.Author, routes.CreateFlag)
	// api.Patch("/flag", auth.Author, routes.UpdateFlag)
	api.Delete("/flag", auth.Author, routes.DeleteFlag)

	api.Patch("/config", auth.Admin, routes.UpdateConfig)

	if log.GetLevel() == log.DebugLevel {
		// Serve frontend in development mode
		frontendAddr := "http://127.0.0.1:5173/"
		if os.Getenv("FRONTEND_ADDR") != "" {
			frontendAddr = os.Getenv("FRONTEND_ADDR")
		}
		app.Use(func(c *fiber.Ctx) error {
			err := proxy.Do(c, frontendAddr+c.Path()[1:])
			if err != nil {
				return err
			}
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		})
	} else {
		// 404 handler
		app.Use(func(c *fiber.Ctx) error {
			return utils.Error(c, fiber.StatusNotFound, consts.EndpointNotFound)
		})
	}

	return app
}
