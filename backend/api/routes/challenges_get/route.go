package challenges_get

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	role := c.Locals("role").(db.UserRole)

	all := utils.In(role, []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	challenges, err := GetChallenges(c.Context(), uid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}

	return c.Status(fiber.StatusOK).JSON(challenges)
}
