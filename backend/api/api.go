package api

import (
	"encoding/json"
	"os"
	"trxd/api/auth"
	"trxd/api/routes/category_create"
	"trxd/api/routes/category_delete"
	"trxd/api/routes/challenge_create"
	"trxd/api/routes/challenge_delete"
	"trxd/api/routes/challenge_flag_create"
	"trxd/api/routes/challenge_flag_delete"
	"trxd/api/routes/challenge_get"
	"trxd/api/routes/challenge_submit"
	"trxd/api/routes/challenges_get"
	"trxd/api/routes/config_update"
	"trxd/api/routes/team_get"
	"trxd/api/routes/team_join"
	"trxd/api/routes/team_password"
	"trxd/api/routes/team_register"
	"trxd/api/routes/team_update"
	"trxd/api/routes/teams_get"
	"trxd/api/routes/user_get"
	"trxd/api/routes/user_info"
	"trxd/api/routes/user_login"
	"trxd/api/routes/user_logout"
	"trxd/api/routes/user_password"
	"trxd/api/routes/user_register"
	"trxd/api/routes/user_update"
	"trxd/api/routes/users_get"
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

	// TODO: make consts for routes
	api.Post("/register", auth.NoAuth, user_register.Route)
	api.Post("/login", auth.NoAuth, user_login.Route)
	api.Post("/logout", auth.NoAuth, user_logout.Route)

	api.Get("/info", auth.Spectator, user_info.Route)
	api.Get("/challenges", auth.Spectator, auth.Team, challenges_get.Route)
	api.Get("/challenges/:id", auth.Spectator, auth.Team, challenge_get.Route)

	// api.Get("/scoreboard", auth.NoAuth)

	api.Get("/users", auth.NoAuth, users_get.Route)
	api.Get("/users/:id", auth.NoAuth, user_get.Route)
	api.Patch("/users", auth.Player, user_update.Route)
	api.Patch("/users/password", auth.Admin, user_password.Route)

	api.Get("/teams", auth.NoAuth, teams_get.Route)
	api.Get("/teams/:id", auth.NoAuth, team_get.Route)
	api.Post("/teams", auth.Player, team_register.Route)
	api.Put("/teams", auth.Player, team_join.Route)
	// api.Get("/teams/join/:token", routes.JoinTeamWithToken)
	api.Patch("/teams", auth.Player, auth.Team, team_update.Route)
	api.Patch("/teams/password", auth.Admin, team_password.Route)

	api.Post("/submit", auth.Player, auth.Team, challenge_submit.Route)

	// api.Post("/instance", auth.Player, routes.CreateInstance)
	// api.Patch("/instance", auth.Player, routes.ExtendInstance)
	// api.Delete("/instance", auth.Player, routes.DeleteInstance)

	api.Post("/category", auth.Author, category_create.Route)
	// api.Patch("/category", auth.Author, routes.UpdateCategory)
	api.Delete("/category", auth.Author, category_delete.Route)

	api.Post("/challenges", auth.Author, challenge_create.Route)
	// api.Patch("/challenges", routes.UpdateChallenge)
	api.Delete("/challenges", auth.Author, challenge_delete.Route)

	api.Post("/flag", auth.Author, challenge_flag_create.Route)
	// api.Patch("/flag", auth.Author, routes.UpdateFlag)
	api.Delete("/flag", auth.Author, challenge_flag_delete.Route)

	api.Patch("/config", auth.Admin, config_update.Route)

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
