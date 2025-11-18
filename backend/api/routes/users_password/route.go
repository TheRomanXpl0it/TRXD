package users_password

import (
	"trxd/api/validator"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		UserID *int32 `json:"user_id" validate:"required,id"`
		// TODO: password new only for ourself
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	newPassword, err := crypto_utils.GeneratePassword()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorGeneratingPassword, err)
	}

	err = ResetUserPassword(c.Context(), *data.UserID, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingUserPassword, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"new_password": newPassword,
	})
}
