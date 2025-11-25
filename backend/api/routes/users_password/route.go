package users_password

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
		UserID      *int32 `json:"user_id" validate:"omitnil,id"`
		NewPassword string `json:"new_password" validate:"omitempty,password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	var uid int32
	role := c.Locals("role").(sqlc.UserRole)
	if role == sqlc.UserRoleAdmin && data.UserID != nil {
		uid = *data.UserID
	} else {
		uid = c.Locals("uid").(int32)
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

	err = ResetUserPassword(c.Context(), uid, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingUserPassword, err)
	}

	if data.NewPassword != "" {
		return c.SendStatus(fiber.StatusOK)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"new_password": newPassword,
	})
}
