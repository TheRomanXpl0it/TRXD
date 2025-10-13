package api

import (
	"context"
	"fmt"
	"time"
	"trxd/api/middlewares"
	"trxd/api/routes/categories_create"
	"trxd/api/routes/categories_delete"
	"trxd/api/routes/categories_update"
	"trxd/api/routes/challenges_all_get"
	"trxd/api/routes/challenges_create"
	"trxd/api/routes/challenges_delete"
	"trxd/api/routes/challenges_get"
	"trxd/api/routes/challenges_update"
	"trxd/api/routes/configs_get"
	"trxd/api/routes/configs_update"
	"trxd/api/routes/flags_create"
	"trxd/api/routes/flags_delete"
	"trxd/api/routes/flags_update"
	"trxd/api/routes/instances_create"
	"trxd/api/routes/instances_delete"
	"trxd/api/routes/instances_update"
	"trxd/api/routes/submissions_create"
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
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/tde-nico/log"
)

var (
	noAuth    = middlewares.NoAuth
	spectator = middlewares.Spectator
	player    = middlewares.Player
	team      = middlewares.Team
	author    = middlewares.Author
	admin     = middlewares.Admin
)

func SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: consts.Name,
	})

	SetupFeatures(app)
	SetupApi(app)

	app.Static("/", "./frontend")
	app.Static("/static", "./static")

	app.Use("/attachments", spectator, team, middlewares.Attachments)
	app.Static("/attachments", "./attachments", fiber.Static{
		Download: true,
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return utils.Error(c, fiber.StatusNotFound, consts.NotFound)
	})

	return app
}

func SetupFeatures(app *fiber.App) {
	if !consts.Testing {
		app.Use(func(c *fiber.Ctx) error {
			defer func() {
				r := recover()
				if r == nil {
					return
				}
				log.Critical("Panic recovered:", "crit", r)
				_ = utils.Error(c, fiber.StatusInternalServerError, consts.InternalServerError)
			}()
			return c.Next()
		})

		// app.Use(limiter.New())
	}

	app.Use(compress.New())

	app.Use(csrf.New(csrf.Config{
		CookieSameSite:    fiber.CookieSameSiteLaxMode,
		CookieSessionOnly: true,
		Expiration:        1 * time.Hour,
		Session:           db.Store,
	}))

	app.Use(favicon.New(favicon.Config{
		File: "./static/favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Get("/monitor", admin, monitor.New(monitor.Config{
		Title: consts.Name + " Monitor",
	}))
}

func SetupApi(app *fiber.App) {
	var api fiber.Router
	if log.GetLevel() == log.DebugLevel {
		api = app.Group("/api", middlewares.Debug)
	} else {
		api = app.Group("/api")
	}

	api.Post("/register", noAuth, users_register.Route)
	api.Post("/login", noAuth, users_login.Route)
	api.Post("/logout", noAuth, users_logout.Route)
	api.Get("/info", noAuth, users_info.Route)
	api.Get("/scoreboard", noAuth, teams_scoreboard.Route)

	api.Patch("/users", player, users_update.Route)
	api.Patch("/users/password", admin, users_password.Route)
	api.Get("/users", noAuth, users_all_get.Route)
	api.Get("/users/:id", noAuth, users_get.Route)

	mode, err := db.GetConfig(context.Background(), "user-mode")
	if err != nil {
		log.Error("Failed to get user-mode config:", "err", err)
		mode = fmt.Sprint(consts.DefaultConfigs["user-mode"])
	}
	if mode != "true" {
		api.Post("/teams/register", player, teams_register.Route)
		api.Post("/teams/join", player, teams_join.Route)
		// api.Get("/teams/join/:token", player, teams_join_token.Route)
		api.Patch("/teams", player, team, teams_update.Route)
		api.Patch("/teams/password", admin, teams_password.Route)
		api.Get("/teams", noAuth, teams_all_get.Route)
		api.Get("/teams/:id", noAuth, teams_get.Route)
	}

	api.Post("/categories", author, categories_create.Route)
	api.Patch("/categories", author, categories_update.Route)
	api.Delete("/categories", author, categories_delete.Route)

	api.Post("/challenges", author, challenges_create.Route)
	api.Patch("/challenges", author, challenges_update.Route)
	api.Delete("/challenges", author, challenges_delete.Route)
	api.Get("/challenges", spectator, team, challenges_all_get.Route)
	api.Get("/challenges/:id", spectator, team, challenges_get.Route)

	api.Post("/instances", player, team, instances_create.Route)
	api.Patch("/instances", player, team, instances_update.Route)
	api.Delete("/instances", player, team, instances_delete.Route)
	// api.Get("/instances", admin, instances_get.Route)

	api.Post("/submissions", spectator, team, submissions_create.Route)
	// api.Get("/submissions", admin, submissions_get.Route)
	// api.Delete("/submissions", admin, submissions_delete.Route)

	api.Post("/tags", author, tags_create.Route)
	api.Patch("/tags", author, tags_update.Route)
	api.Delete("/tags", author, tags_delete.Route)

	api.Post("/flags", author, flags_create.Route)
	api.Patch("/flags", author, flags_update.Route)
	api.Delete("/flags", author, flags_delete.Route)

	api.Get("/configs", admin, configs_get.Route)
	api.Patch("/configs", admin, configs_update.Route)
}
