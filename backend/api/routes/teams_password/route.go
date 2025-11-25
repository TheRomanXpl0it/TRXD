package teams_password

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		TeamID      *int32 `json:"team_id" validate:"omitnil,id"`
		NewPassword string `json:"new_password" validate:"omitempty,password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	var tid int32
	role := c.Locals("role").(sqlc.UserRole)
	if role == sqlc.UserRoleAdmin && data.TeamID != nil {
		tid = *data.TeamID
	} else {
		tid = c.Locals("tid").(int32)
	}

	if tid < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidTeamID)
	}

	var newPassword string
	if data.NewPassword != "" {
		newPassword = data.NewPassword
	} else {
		newPassword, err = crypto_utils.GeneratePassword()
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorGeneratingPassword, err)
		}
	}

	err = ResetTeamPassword(c.Context(), tid, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingTeamPassword, err)
	}

	if data.NewPassword != "" {
		return c.SendStatus(fiber.StatusOK)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"new_password": newPassword,
	})
}
