package users_role

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		UserID  *int32        `json:"user_id" validate:"required,id"`
		NewRole sqlc.UserRole `json:"new_role" validate:"required,user_role"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	if data.NewRole == sqlc.UserRoleAdmin {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidRole)
	}

	err = ChangeUserRole(c.Context(), *data.UserID, data.NewRole)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorChangingUserRole, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
