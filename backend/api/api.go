package api

import (
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
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tde-nico/log"
)

func SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: consts.Name,
	})

	// app.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestSpeed, // 1
	// }))

	// app.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }))

	// app.Use(limiter.New(limiter.Config{
	// 	Next: func(c *fiber.Ctx) bool {
	// 		return c.IP() == "127.0.0.1"
	// 	},
	// 	Max:        20,
	// 	Expiration: 30 * time.Second,
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		return c.Get("x-forwarded-for")
	// 	},
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		return c.SendFile("./toofast.html")
	// 	},
	// 	Storage: myCustomStorage{},
	// }))

	// app.Use(favicon.New(favicon.Config{
	// 	File: "./favicon.ico",
	// 	URL:  "/favicon.ico",
	// }))

	// app.Use(recover.New())

	// app.Get("/foo/:sleepTime", timeout.New(h, 2*time.Second))
	// app.Get("/foo", timeout.NewWithContext(handler, 10*time.Second))

	// app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:1337", // TODO: edit
		AllowCredentials: true,
	}))

	var api fiber.Router
	if log.GetLevel() == log.DebugLevel {
		api = app.Group("/api", middlewares.Debug)
	} else {
		api = app.Group("/api")
	}
	noAuth := middlewares.NoAuth
	spectator := middlewares.Spectator
	player := middlewares.Player
	team := middlewares.Team
	author := middlewares.Author
	admin := middlewares.Admin

	// TODO: make this resource static: countries.json
	api.Get("/countries", func(c *fiber.Ctx) error { return c.JSON(consts.Countries) })

	api.Post("/register", noAuth, users_register.Route)
	api.Post("/login", noAuth, users_login.Route)
	api.Post("/logout", noAuth, users_logout.Route)
	api.Get("/info", spectator, users_info.Route)
	api.Get("/scoreboard", noAuth, teams_scoreboard.Route)

	api.Patch("/users", player, users_update.Route)
	api.Patch("/users/password", admin, users_password.Route)
	api.Get("/users", noAuth, users_all_get.Route)
	api.Get("/users/:id", noAuth, users_get.Route)

	api.Post("/teams/register", player, teams_register.Route)
	api.Post("/teams/join", player, teams_join.Route)
	// api.Get("/teams/join/:token", player, teams_join_token.Route)
	api.Patch("/teams", player, team, teams_update.Route)
	api.Patch("/teams/password", admin, teams_password.Route)
	api.Get("/teams", noAuth, teams_all_get.Route)
	api.Get("/teams/:id", noAuth, teams_get.Route)

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

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return utils.Error(c, fiber.StatusNotFound, consts.EndpointNotFound)
	})

	return app
}
