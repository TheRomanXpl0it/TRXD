package challenges_all_get

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/pluginsapi"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	tid := c.Locals("tid").(int32)
	role := c.Locals("role").(sqlc.UserRole)
	
	payload := struct {
		uid int32
		tid int32
		role sqlc.UserRole
	}{
		uid: uid,
		tid: tid,
		role: role,
	}
	payload,err := pluginsApi.DispatchEvent(c.Context(),"challengesGet",payload)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorPassingDataToPlugins, err)
	}

	
	

	all := utils.In(role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	challenges, err := GetChallenges(c.Context(), uid, tid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}
	
	challenges,err = pluginsApi.DispatchEvent(c.Context(),"challengesGet",challenges)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorPassingDataToPlugins, err)
	}

	return c.Status(fiber.StatusOK).JSON(challenges)
}
