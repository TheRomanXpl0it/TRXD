package middlewares

import (
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func timeGate(c *fiber.Ctx, key string, check func(now, t time.Time) bool, errMsg string) error {
	start, err := db.GetConfig(c.Context(), key)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if start == "" {
		return c.Next()
	}

	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorParsingTime, err)
	}
	if check(time.Now(), startTime) {
		return c.Next()
	}

	role := c.Locals("role")
	if role == nil ||
		(role.(sqlc.UserRole) != sqlc.UserRoleAuthor && role.(sqlc.UserRole) != sqlc.UserRoleAdmin) {
		return utils.Error(c, fiber.StatusForbidden, errMsg)
	}

	return c.Next()
}

func Start(c *fiber.Ctx) error {
	return timeGate(
		c,
		"start-time",
		func(now, t time.Time) bool { return !now.Before(t) },
		consts.NotStartedYet,
	)
}

func End(c *fiber.Ctx) error {
	return timeGate(
		c,
		"end-time",
		func(now, t time.Time) bool { return !now.After(t) },
		consts.AlreadyEnded,
	)
}
