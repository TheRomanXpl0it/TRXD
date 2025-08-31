package teams_password

import (
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		TeamID *int32 `json:"team_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.TeamID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.TeamID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidTeamID)
	}

	newPassword, err := crypto_utils.GeneratePassword()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorGeneratingPassword, err)
	}

	err = ResetTeamPassword(c.Context(), *data.TeamID, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingTeamPassword, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"new_password": newPassword,
	})
}
