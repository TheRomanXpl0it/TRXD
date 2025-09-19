package instances_create

import (
	"fmt"
	"time"
	"trxd/api/validator"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/instancer"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role").(sqlc.UserRole)
	tid := c.Locals("tid").(int32)
	if tid == -1 {
		return utils.Error(c, fiber.StatusForbidden, consts.TeamNotFound)
	}

	var data struct {
		ChallID *int32 `json:"chall_id" validate:"required,challenge_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	secret, err := db.GetConfig(c.Context(), "secret")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if secret == nil || secret.Value == "" {
		return utils.Error(c, fiber.StatusForbidden, consts.DisabledInstances)
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
	if instance != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyAnActiveInstance)
	}

	if chall.DockerConfig.Lifetime == 0 {
		return utils.Error(c, fiber.StatusInternalServerError, consts.MissingLifetime, fmt.Errorf(consts.MissingLifetime))
	}
	lifetime := time.Second * time.Duration(chall.DockerConfig.Lifetime.(int64))
	expires_at := time.Now().Add(lifetime)
	var internalPort *int32
	if chall.Info.Port != 0 {
		internalPort = &chall.Info.Port
	}

	host, port, err := instancer.CreateInstance(c.Context(), tid, *data.ChallID, internalPort, expires_at, chall.Info.Type, chall.DockerConfig)
	if err != nil {
		switch err.Error() {
		case "[race condition]":
			return utils.Error(c, fiber.StatusConflict, consts.AlreadyAnActiveInstance)
		case "[no image or compose]":
			return utils.Error(c, fiber.StatusBadRequest, consts.InvalidImage)
		default:
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingInstance, err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		Host    string `json:"host"`
		Port    *int32 `json:"port,omitempty"`
		Timeout int    `json:"timeout"`
	}{
		Host:    host,
		Port:    port,
		Timeout: max(int(time.Until(expires_at).Seconds()), 0),
	})
}
