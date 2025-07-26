package api

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/tde-nico/log"
)

const (
	minPasswordLength      = 8
	maxPasswordLength      = 64
	maxNameLength          = 64
	maxEmailLength         = 256
	invalidJSON            = "Invalid JSON format"
	missingRequiredFields  = "Missing required fields"
	shortPassword          = "Password must be at least 8 characters long"
	longPassword           = "Password must not exceed 64 characters"
	longName               = "Username must not exceed 64 characters"
	longEmail              = "Email must not exceed 256 characters"
	invalidEmail           = "Invalid email format"
	userAlreadyExists      = "User already exists"
	errorRegisteringUser   = "Error registering user"
	errorFetchingSession   = "Error fetching session"
	errorSavingSession     = "Error saving session"
	errorDestroyingSession = "Error destroying session"
	errorLoggingIn         = "Error logging in"
	invalidCredentials     = "Invalid email or password"
	unauthorized           = "Unauthorized"
	alreadyInTeam          = "Already in a team"
	errorRegisteringTeam   = "Error registering team"
	teamAlreadyExists      = "Team already exists"
)

var UserRegex = regexp.MustCompile(`(^[^@\s]+@[^@\s]+\.[^@\s]+$)`)

func apiError(c *fiber.Ctx, status int, message string, err ...error) error {
	if err != nil {
		log.Error("API Error:", "desc", message, "err", err)
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// TODO: set store as redis + set configs (expire time, etc.)
var store = session.New(session.Config{})

func AuthRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return apiError(c, fiber.StatusUnauthorized, unauthorized)
	}
	c.Locals("uid", uid)

	return c.Next()
}

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

	// gob.Register(&uuid.UUID{})

	app.Post("/register", register)
	app.Post("/login", login)
	app.Post("/logout", logout)
	//app.Post("/register-team", AuthRequired, registerTeam)

	app.Get("/test", AuthRequired, func(c *fiber.Ctx) error {
		uid := c.Locals("uid")
		return c.JSON(fiber.Map{
			"uid": uid,
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return app
}
