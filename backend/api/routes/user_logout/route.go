package user_logout

import (
	"trxd/api/middlewares"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid")
	if uid == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.NotLoggedIn)
	}

	sess, err := middlewares.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	err = sess.Destroy()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDestroyingSession, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
