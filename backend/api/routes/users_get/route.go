package users_get

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role")

	allData := false
	if role != nil {
		allData = utils.In(role.(db.UserRole), []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	}
	usersData, err := GetUsers(c.Context(), allData)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(usersData)
}
