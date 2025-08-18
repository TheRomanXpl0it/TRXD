package api

import (
	"os"
	"trxd/api/middlewares"
	"trxd/api/routes/category_create"
	"trxd/api/routes/category_delete"
	"trxd/api/routes/category_update"
	"trxd/api/routes/challenge_create"
	"trxd/api/routes/challenge_delete"
	"trxd/api/routes/challenge_flag_create"
	"trxd/api/routes/challenge_flag_delete"
	"trxd/api/routes/challenge_flag_update"
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
	"trxd/api/routes/teams_scoreboard"
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
		api = app.Group("/api", middlewares.Debug)
	} else {
		api = app.Group("/api")
	}

	// TODO: make this resource static: countries.json
	api.Get("/countries", func(c *fiber.Ctx) error { return c.JSON(consts.Countries) })

	api.Post("/register", middlewares.NoAuth, user_register.Route)
	api.Post("/login", middlewares.NoAuth, user_login.Route)
	api.Post("/logout", middlewares.NoAuth, user_logout.Route)

	api.Get("/info", middlewares.Spectator, user_info.Route)
	api.Get("/challenges", middlewares.Spectator, middlewares.Team, challenges_get.Route)
	api.Get("/challenges/:id", middlewares.Spectator, middlewares.Team, challenge_get.Route)

	api.Get("/users", middlewares.NoAuth, users_get.Route)
	api.Get("/users/:id", middlewares.NoAuth, user_get.Route)
	api.Patch("/users", middlewares.Player, user_update.Route)
	api.Patch("/users/password", middlewares.Admin, user_password.Route)

	api.Get("/scoreboard", middlewares.NoAuth, teams_scoreboard.Route)
	api.Get("/teams", middlewares.NoAuth, teams_get.Route)
	api.Get("/teams/:id", middlewares.NoAuth, team_get.Route)
	api.Post("/teams", middlewares.Player, team_register.Route)
	api.Put("/teams", middlewares.Player, team_join.Route)
	// api.Get("/teams/join/:token", routes.JoinTeamWithToken)
	api.Patch("/teams", middlewares.Player, middlewares.Team, team_update.Route)
	api.Patch("/teams/password", middlewares.Admin, team_password.Route)

	api.Post("/submit", middlewares.Player, middlewares.Team, challenge_submit.Route)

	// api.Post("/instance", middlewares.Player, routes.CreateInstance)
	// api.Patch("/instance", middlewares.Player, routes.ExtendInstance)
	// api.Delete("/instance", middlewares.Player, routes.DeleteInstance)

	api.Post("/category", middlewares.Author, category_create.Route)
	api.Patch("/category", middlewares.Author, category_update.Route)
	api.Delete("/category", middlewares.Author, category_delete.Route)

	api.Post("/challenges", middlewares.Author, challenge_create.Route)
	// api.Patch("/challenges", routes.UpdateChallenge)
	api.Delete("/challenges", middlewares.Author, challenge_delete.Route)

	api.Post("/flag", middlewares.Author, challenge_flag_create.Route)
	api.Patch("/flag", middlewares.Author, challenge_flag_update.Route)
	api.Delete("/flag", middlewares.Author, challenge_flag_delete.Route)

	api.Patch("/config", middlewares.Admin, config_update.Route)

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
