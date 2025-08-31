package users_password

import (
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		UserID *int32 `json:"user_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.UserID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.UserID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidUserID)
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
