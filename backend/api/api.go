package api

import (
	"trxd/api/auth"
	"trxd/api/routes"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "TRXd",
	})

	// TODO: fix CORS settings
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:5173",
	// 	AllowCredentials: true,
	// 	// AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	// 	// AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	// 	// ExposeHeaders:    "Content-Length,Content-Type",
	// }))

	app.Post("/register", routes.Register)
	app.Post("/login", routes.Login)
	app.Post("/logout", routes.Logout)

	app.Post("/register-team", auth.PlayerRequired, routes.RegisterTeam)
	app.Post("/join-team", auth.PlayerRequired, routes.JoinTeam)
	app.Post("/submit", auth.PlayerRequired, routes.Submit)

	app.Post("/create-category", auth.AuthorRequired, routes.CreateCategory)
	app.Post("/create-challenge", auth.AuthorRequired, routes.CreateChallenge)
	app.Post("/create-flag", auth.AuthorRequired, routes.CreateFlag)

	app.Post("/update-config", auth.AdminRequired, routes.UpdateConfig)

	// TODO: remove this endpoint
	//! ############################## TEST ENDPOINT ##############################
	app.Get("/test", auth.AuthRequired, func(c *fiber.Ctx) error {
		uid := c.Locals("uid")
		team, err := db.GetTeamFromUser(c.Context(), uid.(int32))
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
		}

		return c.JSON(fiber.Map{
			"uid":  uid,
			"team": team,
		})
	})
	//! ############################## TEST ENDPOINT ##############################

	app.Use(func(c *fiber.Ctx) error {
		return utils.Error(c, fiber.StatusNotFound, consts.EndpointNotFound)
	})

	return app
}
