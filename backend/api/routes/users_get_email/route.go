package users_get_email

import (
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	userEmail := c.Query("email")
	if userEmail == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidEmail)
	}

	valid, err := validator.Var(c, userEmail, "user_email")
	if err != nil || !valid {
		return err
	}

	userData, err := GetUserByEmail(c.Context(), userEmail)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if userData == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.UserNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(userData)
}
