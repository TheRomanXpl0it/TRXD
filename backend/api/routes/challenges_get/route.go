package challenges_get

import (
	"trxd/api/validator"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	tid := c.Locals("tid").(int32)
	role := c.Locals("role").(sqlc.UserRole)

	challengeIDInt, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallengeID)
	}
	challengeID := int32(challengeIDInt)
	valid, err := validator.Var(c, challengeID, "id")
	if err != nil || !valid {
		return err
	}

	all := utils.In(role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	challenge, err := GetChallenge(c.Context(), challengeID, uid, tid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(challenge)
}
