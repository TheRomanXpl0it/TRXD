package users_get

import (
	"trxd/api/validator"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid")
	role := c.Locals("role")

	userID, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidUserID)
	}
	valid, err := validator.Var(c, userID, "user_id")
	if err != nil || !valid {
		return err
	}

	allData := false
	if uid != nil {
		allData = uid.(int32) == int32(userID)
	}
	if !allData && role != nil {
		allData = utils.In(role.(sqlc.UserRole), []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	}
	userData, err := GetUser(c.Context(), int32(userID), allData)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if userData == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.UserNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(userData)
}
