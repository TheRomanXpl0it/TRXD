package users_all_get

import (
	"math"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role")

	offset := c.QueryInt("offset", 0)
	if offset < 0 || offset > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	limit := c.QueryInt("limit", 0)
	if limit < 0 || limit > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	allData := false
	if role != nil {
		allData = utils.In(role.(sqlc.UserRole), []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	}
	totalUsers, usersData, err := GetUsers(c.Context(), allData, int32(offset), int32(limit))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total": totalUsers,
		"users": usersData,
	})
}
