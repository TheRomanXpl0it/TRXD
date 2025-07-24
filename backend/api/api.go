package api

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	minPasswordLength     = 8
	maxPasswordLength     = 64
	maxUsernameLength     = 64
	maxEmailLength        = 256
	invalidJSON           = "Invalid JSON format"
	missingRequiredFields = "Missing required fields"
	shortPassword         = "Password must be at least 8 characters long"
	longPassword          = "Password must not exceed 64 characters"
	longUsername          = "Username must not exceed 64 characters"
	longEmail             = "Email must not exceed 256 characters"
	invalidEmail          = "Invalid email format"
	userAlreadyExists     = "User already exists"
	errorRegisteringUser  = "Error registering user"
	errorCreatingSession  = "Error creating session"
)

var UserRegex = regexp.MustCompile(`(^[^@\s]+@[^@\s]+\.[^@\s]+$)`)

func apiError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// TODO: set store as redis + set configs (expire time, etc.)
var store = session.New()

func SetupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "TRXd",
	})

	app.Post("/register", registerPost)
	app.Post("/login", loginPost)
	app.Post("/logout", logoutPost)

	return app
}
