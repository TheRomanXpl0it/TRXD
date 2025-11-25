package teams_register

import (
	"strings"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name" validate:"required,team_name"`
		Password string `json:"password" validate:"required,password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	uid := c.Locals("uid").(int32)

	team, err := db.GetTeamFromUser(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if team != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyInTeam)
	}

	tx, err := db.BeginTx(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorBeginningTransaction, err)
	}
	defer tx.Rollback()

	team, err = RegisterTeam(c.Context(), tx, data.Name, data.Password, uid)
	if err != nil {
		if strings.HasPrefix(err.Error(), "[race condition]") {
			return utils.Error(c, fiber.StatusConflict, consts.AlreadyInTeam)
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
	}
	if team == nil {
		return utils.Error(c, fiber.StatusConflict, consts.TeamAlreadyExists)
	}

	err = tx.Commit()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCommittingTransaction, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
