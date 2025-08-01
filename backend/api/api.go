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

	app.Patch("/user", auth.PlayerRequired, routes.UpdateUser)
	app.Post("/register-team", auth.PlayerRequired, routes.RegisterTeam)
	app.Post("/join-team", auth.PlayerRequired, routes.JoinTeam)
	app.Patch("/team", auth.PlayerRequired, routes.UpdateTeam)
	app.Post("/submit", auth.PlayerRequired, routes.Submit)
	// app.Post("/instance", auth.PlayerRequired, routes.CreateInstance)
	// app.Patch("/instance", auth.PlayerRequired, routes.ExtendInstance)
	// app.Delete("/instance", auth.PlayerRequired, routes.DeleteInstance)

	app.Post("/category", auth.AuthorRequired, routes.CreateCategory)
	// app.Patch("/category", auth.AuthorRequired, routes.UpdateCategory)
	app.Delete("/category", auth.AuthorRequired, routes.DeleteCategory)
	app.Post("/challenge", auth.AuthorRequired, routes.CreateChallenge)
	// app.Patch("/challenge", auth.AuthorRequired, routes.UpdateChallenge)
	app.Delete("/challenge", auth.AuthorRequired, routes.DeleteChallenge)
	app.Post("/flag", auth.AuthorRequired, routes.CreateFlag)
	// app.Patch("/flag", auth.AuthorRequired, routes.UpdateFlag)
	app.Delete("/flag", auth.AuthorRequired, routes.DeleteFlag)

	app.Patch("/config", auth.AdminRequired, routes.UpdateConfig)
	app.Post("/reset-user-password", auth.AdminRequired, routes.ResetUserPassword)
	app.Post("/reset-team-password", auth.AdminRequired, routes.ResetTeamPassword)

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
