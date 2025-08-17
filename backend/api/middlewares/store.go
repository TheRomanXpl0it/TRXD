package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// TODO: set store as redis
var Store = session.New(session.Config{
	Expiration:     30 * 24 * time.Hour,
	CookiePath:     "/",
	CookieSameSite: fiber.CookieSameSiteLaxMode,
})
