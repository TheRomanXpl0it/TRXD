package instances_delete

import (
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/instancer"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role").(sqlc.UserRole)
	tid := c.Locals("tid").(int32)
	if tid == -1 {
		return utils.Error(c, fiber.StatusForbidden, consts.TeamNotFound)
	}

	var data struct {
		ChallID *int32 `json:"chall_id" validate:"required,id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	chall, err := db.GetChallenge(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if chall == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	if chall.Info.Hidden && !utils.In(role,
		[]sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin}) {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}
	if chall.Info.Type == sqlc.DeployTypeNormal {
		return utils.Error(c, fiber.StatusBadRequest, consts.ChallengeNotInstanciable)
	}

	instance, err := instancer.GetInstance(c.Context(), *data.ChallID, tid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingInstance, err)
	}
	if instance == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.InstanceNotFound)
	}

	err = instancer.DeleteInstance(c.Context(), tid, *data.ChallID, instance.DockerID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingInstance, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
