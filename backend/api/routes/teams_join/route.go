package teams_join

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) < consts.MinPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.ShortPassword)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Name) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}

	uid := c.Locals("uid").(int32)

	team, err := db.GetTeamFromUser(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if team != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyInTeam)
	}

	team, err = JoinTeam(c.Context(), data.Name, data.Password, uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
	}
	if team == nil {
		return utils.Error(c, fiber.StatusConflict, consts.InvalidTeamCredentials)
	}

	return c.SendStatus(fiber.StatusOK)
}
