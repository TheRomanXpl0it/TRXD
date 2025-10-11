package challenges_all_get

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	tid := c.Locals("tid").(int32)
	role := c.Locals("role").(sqlc.UserRole)

	all := utils.In(role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	challenges, err := GetChallenges(c.Context(), uid, tid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}

	return c.Status(fiber.StatusOK).JSON(challenges)
}
