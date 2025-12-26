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

	start := c.QueryInt("start", 0)
	if start < 0 || start > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	end := c.QueryInt("end", 0)
	if end < 0 || end > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	if end != 0 && end <= start {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	allData := false
	if role != nil {
		allData = utils.In(role.(sqlc.UserRole), []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	}
	usersData, err := GetUsers(c.Context(), allData, int32(start), int32(end))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(usersData)
}
