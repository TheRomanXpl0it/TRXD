package api

import (
	"regexp"
	"trxd/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/tde-nico/log"
)

const (
	MinPasswordLength      = 8
	MaxPasswordLength      = 64
	MaxNameLength          = 64
	MaxEmailLength         = 256
	InvalidJSON            = "Invalid JSON format"
	MissingRequiredFields  = "Missing required fields"
	ShortPassword          = "Password must be at least 8 characters long"
	LongPassword           = "Password must not exceed 64 characters"
	LongName               = "Username must not exceed 64 characters"
	LongEmail              = "Email must not exceed 256 characters"
	InvalidEmail           = "Invalid email format"
	UserAlreadyExists      = "User already exists"
	ErrorRegisteringUser   = "Error registering user"
	ErrorFetchingSession   = "Error fetching session"
	ErrorSavingSession     = "Error saving session"
	ErrorDestroyingSession = "Error destroying session"
	ErrorLoggingIn         = "Error logging in"
	InvalidCredentials     = "Invalid email or password"
	Unauthorized           = "Unauthorized"
	AlreadyInTeam          = "Already in a team"
	ErrorRegisteringTeam   = "Error registering team"
	TeamAlreadyExists      = "Team already exists"
	ErrorFetchingTeam      = "Error fetching team"
	InvalidTeamCredentials = "Invalid name or password"
)

var UserRegex = regexp.MustCompile(`(^[^@\s]+@[^@\s]+\.[^@\s]+$)`)

// TODO: set store as redis + set configs (expire time, etc.)
var store = session.New(session.Config{})

func apiError(c *fiber.Ctx, status int, message string, err ...error) error {
	if err != nil {
		log.Error("API Error:", "desc", message, "err", err)
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}

func AuthRequired(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return apiError(c, fiber.StatusUnauthorized, Unauthorized)
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

	app.Post("/register", register)
	app.Post("/login", login)
	app.Post("/logout", logout)
	app.Post("/register-team", AuthRequired, registerTeam)
	app.Post("/join-team", AuthRequired, joinTeam)

	// TODO: remove this endpoint
	//! ############################## TEST ENDPOINT ##############################
	app.Get("/test", AuthRequired, func(c *fiber.Ctx) error {
		uid := c.Locals("uid")
		team, err := db.GetTeamFromUser(c.Context(), uid.(int32))
		if err != nil {
			return apiError(c, fiber.StatusInternalServerError, ErrorFetchingTeam, err)
		}

		return c.JSON(fiber.Map{
			"uid":  uid,
			"team": team,
		})
	})
	//! ############################## TEST ENDPOINT ##############################

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return app
}
