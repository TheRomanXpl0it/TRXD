package challenges_get

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	role := c.Locals("role").(sqlc.UserRole)

	challengeID, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallengeID)
	}

	all := utils.In(role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	challenge, err := GetChallenge(c.Context(), int32(challengeID), uid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(challenge)
}
