package users_login

import (
	"trxd/api/middlewares"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid")
	if uid != nil {
		return utils.Error(c, fiber.StatusForbidden, consts.AlreadyLoggedIn)
	}

	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Email == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Email) > consts.MaxEmailLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongEmail)
	}

	user, err := LoginUser(c.Context(), data.Email, data.Password)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorLoggingIn, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.InvalidCredentials)
	}

	sess, err := middlewares.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", user.ID)
	err = sess.Save()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
