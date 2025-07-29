package api

import (
	"github.com/gofiber/fiber/v2"
)

func logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorFetchingSession, err)
	}

	err = sess.Destroy()
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorDestroyingSession, err)
	}

	return c.Status(fiber.StatusOK).SendString("")
}
