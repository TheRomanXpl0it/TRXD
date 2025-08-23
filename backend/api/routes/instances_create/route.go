package instances_create

import (
	"time"
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

	chall, err := GetChallenge(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge)
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
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingInstance)
	}
	if instance != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyAnActiveInstance)
	}

	// TODO: GenPort
	port := 1337
	// TODO: GenDomain
	domain := "127.0.0.1"
	if !chall.DockerConfig.Lifetime.Valid {
		// TODO: global lifetime in configs
		chall.DockerConfig.Lifetime.Int32 = 60
		// return utils.Error(c, fiber.StatusInternalServerError, "DockerConfig lifetime should be valid", fmt.Errorf("DockerConfig lifetime should be valid"))
	}
	lifetime := time.Second * time.Duration(chall.DockerConfig.Lifetime.Int32)
	expires_at := time.Now().Add(lifetime)

	raceCondition, err := CreateInstance(c.Context(), tid, *data.ChallID, expires_at, domain, port)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingInstance)
	}
	if raceCondition { // TODO: tests
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyAnActiveInstance)
	}

	// TODO: call the instancer
	log.Info("Creating Instance:", "tid", tid, "challID", *data.ChallID, "expiresAt", expires_at,
		"domain", domain, "port", port, "docker-conf", chall.DockerConfig)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"expires_at": expires_at,
		"domain":     domain,
		"port":       port,
	})
}
