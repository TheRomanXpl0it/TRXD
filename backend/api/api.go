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

func debugMiddleware(c *fiber.Ctx) error {

	reqBody := c.BodyRaw()
	body := string(reqBody)
	var tmp map[string]interface{}
	if err := json.Unmarshal(reqBody, &tmp); err == nil {
		if tmp2, err := json.MarshalIndent(tmp, "", "  "); err == nil {
			body = string(tmp2)
		}
	}
	log.Debug("Request:", "method", c.Method(), "path", c.Path(), "body", body)

	e := c.Next()

	resStatus := c.Response().StatusCode()
	resBody := c.Response().Body()
	body = string(resBody)
	tmp = map[string]interface{}{}
	if err := json.Unmarshal(resBody, &tmp); err == nil {
		if tmp2, err := json.MarshalIndent(tmp, "", "  "); err == nil {
			body = string(tmp2)
		}
	}
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
	player := api.Group("/player", auth.PlayerRequired)
	author := api.Group("/author", auth.AuthorRequired)
	admin := api.Group("/admin", auth.AdminRequired)

	api.Post("/register", routes.Register)
	api.Post("/login", routes.Login)
	api.Post("/logout", routes.Logout)
	api.Get("/auth", auth.AuthRequired, routes.Auth)

	player.Patch("/user", routes.UpdateUser)
	player.Post("/register-team", routes.RegisterTeam)
	player.Post("/join-team", routes.JoinTeam)
	player.Patch("/team", routes.UpdateTeam)
	player.Post("/submit", routes.Submit)
	// player.Post("/instance", routes.CreateInstance)
	// player.Patch("/instance", routes.ExtendInstance)
	// player.Delete("/instance", routes.DeleteInstance)

	author.Post("/category", routes.CreateCategory)
	// author.Patch("/category", routes.UpdateCategory)
	author.Delete("/category", routes.DeleteCategory)
	author.Post("/challenge", routes.CreateChallenge)
	// author.Patch("/challenge", routes.UpdateChallenge)
	author.Delete("/challenge", routes.DeleteChallenge)
	author.Post("/flag", routes.CreateFlag)
	// author.Patch("/flag", routes.UpdateFlag)
	author.Delete("/flag", routes.DeleteFlag)

	admin.Patch("/config", routes.UpdateConfig)
	admin.Post("/reset-user-password", routes.ResetUserPassword)
	admin.Post("/reset-team-password", routes.ResetTeamPassword)

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
