package routes

import (
	"trxd/api/auth"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	sess, err := auth.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	err = sess.Destroy()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDestroyingSession, err)
	}

	return c.Status(fiber.StatusOK).SendString("")
}
