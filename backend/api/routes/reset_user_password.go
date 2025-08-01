package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func ResetUserPassword(c *fiber.Ctx) error {
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

	newPassword, err := utils.GenerateRandPass()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorGeneratingPassword, err)
	}

	err = db.ResetUserPassword(c.Context(), *data.UserID, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingUserPassword, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"new_password": newPassword})
}
