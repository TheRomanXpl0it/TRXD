package api

import (
	"os"
	"trxd/api/middlewares"
	"trxd/api/routes/categories_create"
	"trxd/api/routes/categories_delete"
	"trxd/api/routes/categories_update"
	"trxd/api/routes/challenges_all_get"
	"trxd/api/routes/challenges_create"
	"trxd/api/routes/challenges_delete"
	"trxd/api/routes/challenges_get"
	"trxd/api/routes/challenges_update"
	"trxd/api/routes/configs_update"
	"trxd/api/routes/flags_create"
	"trxd/api/routes/flags_delete"
	"trxd/api/routes/flags_submit"
	"trxd/api/routes/flags_update"
	"trxd/api/routes/tags_create"
	"trxd/api/routes/tags_delete"
	"trxd/api/routes/tags_update"
	"trxd/api/routes/teams_all_get"
	"trxd/api/routes/teams_get"
	"trxd/api/routes/teams_join"
	"trxd/api/routes/teams_password"
	"trxd/api/routes/teams_register"
	"trxd/api/routes/teams_scoreboard"
	"trxd/api/routes/teams_update"
	"trxd/api/routes/users_all_get"
	"trxd/api/routes/users_get"
	"trxd/api/routes/users_info"
	"trxd/api/routes/users_login"
	"trxd/api/routes/users_logout"
	"trxd/api/routes/users_password"
	"trxd/api/routes/users_register"
	"trxd/api/routes/users_update"
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

	api.Post("/users/register", middlewares.NoAuth, users_register.Route)
	api.Post("/users/login", middlewares.NoAuth, users_login.Route)
	api.Post("/users/logout", middlewares.NoAuth, users_logout.Route)
	api.Patch("/users/update", middlewares.Player, users_update.Route)
	api.Patch("/users/password", middlewares.Admin, users_password.Route)
	api.Get("/users/info", middlewares.Spectator, users_info.Route)
	api.Get("/users", middlewares.NoAuth, users_all_get.Route)
	api.Get("/users/:id", middlewares.NoAuth, users_get.Route)

	api.Post("/teams/register", middlewares.Player, teams_register.Route)
	api.Post("/teams/join", middlewares.Player, teams_join.Route)
	// api.Get("/teams/join/:token", middlewares.Player, teams_join_token.Route)
	api.Patch("/teams/update", middlewares.Player, middlewares.Team, teams_update.Route)
	api.Patch("/teams/password", middlewares.Admin, teams_password.Route)
	api.Get("/teams/scoreboard", middlewares.NoAuth, teams_scoreboard.Route)
	api.Get("/teams", middlewares.NoAuth, teams_all_get.Route)
	api.Get("/teams/:id", middlewares.NoAuth, teams_get.Route)

	api.Post("/categories/create", middlewares.Author, categories_create.Route)
	api.Patch("/categories/update", middlewares.Author, categories_update.Route)
	api.Delete("/categories/delete", middlewares.Author, categories_delete.Route)

	api.Post("/challenges/create", middlewares.Author, challenges_create.Route)
	api.Patch("/challenges/update", middlewares.Author, challenges_update.Route)
	api.Delete("/challenges/delete", middlewares.Author, challenges_delete.Route)
	api.Get("/challenges", middlewares.Spectator, middlewares.Team, challenges_all_get.Route)
	api.Get("/challenges/:id", middlewares.Spectator, middlewares.Team, challenges_get.Route)

	// api.Post("/instances/create", middlewares.Player, instances_create.Route)
	// api.Patch("/instances/update", middlewares.Player, instances_update.Route)
	// api.Delete("/instances/delete", middlewares.Player, instances_delete.Route)

	api.Post("/tags/create", middlewares.Author, tags_create.Route)   // TODO: complete
	api.Patch("/tags/update", middlewares.Author, tags_update.Route)  // TODO: complete
	api.Delete("/tags/delete", middlewares.Author, tags_delete.Route) // TODO: complete

	api.Post("/flags/create", middlewares.Author, flags_create.Route)
	api.Patch("/flags/update", middlewares.Author, flags_update.Route)
	api.Delete("/flags/delete", middlewares.Author, flags_delete.Route)
	api.Post("/flags/submit", middlewares.Player, middlewares.Team, flags_submit.Route)

	api.Patch("/configs/update", middlewares.Admin, configs_update.Route)

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
