package instances_create

import (
	"fmt"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/tde-nico/log"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role").(sqlc.UserRole)
	tid := c.Locals("tid").(int32)
	if tid == -1 {
		return utils.Error(c, fiber.StatusForbidden, consts.TeamNotFound)
	}

	var data struct {
		ChallID *int32 `json:"chall_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.ChallID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallengeID)
	}

	secret, err := db.GetConfig(c.Context(), "secret")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if secret == nil || secret.Value == "" {
		return utils.Error(c, fiber.StatusForbidden, consts.DisabledInstances)
	}

	chall, err := GetChallenge(c.Context(), *data.ChallID)
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

	instance, err := GetInstance(c.Context(), *data.ChallID, tid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingInstance, err)
	}
	if instance != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyAnActiveInstance)
	}

	if !chall.DockerConfig.Lifetime.Valid {
		return utils.Error(c, fiber.StatusInternalServerError, consts.MissingLifetime, fmt.Errorf(consts.MissingLifetime))
	}
	lifetime := time.Second * time.Duration(chall.DockerConfig.Lifetime.Int32)
	expires_at := time.Now().Add(lifetime)

	info, err := CreateInstance(c.Context(), tid, *data.ChallID, expires_at, chall.DockerConfig.HashDomain)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingInstance, err)
	}
	if info == nil { // race condition
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyAnActiveInstance)
	}

	// TODO: call the instancer
	log.Info("Creating Instance:", "tid", tid, "challID", *data.ChallID, "expiresAt", expires_at,
		"host", info.Host, "port", info.Port, "docker-conf", chall.DockerConfig)

	var port *int32
	if !chall.DockerConfig.HashDomain {
		port = &info.Port
	}
	timeout := int(time.Until(expires_at).Seconds())
	if timeout < 0 {
		timeout = 0
	}

	return c.Status(fiber.StatusOK).JSON(struct {
		Host    string `json:"host"`
		Port    *int32 `json:"port,omitempty"`
		Timeout int    `json:"timeout"`
	}{
		Host:    info.Host,
		Port:    port,
		Timeout: timeout,
	})
}
