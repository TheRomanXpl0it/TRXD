package flags_create

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id" validate:"required,id"`
		Flag    string `json:"flag" validate:"required,flag"`
		Regex   bool   `json:"regex"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	flag, err := CreateFlag(c.Context(), *data.ChallID, data.Flag, data.Regex)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingFlag, err)
	}
	if flag == nil {
		return utils.Error(c, fiber.StatusConflict, consts.FlagAlreadyExists)
	}

	return c.SendStatus(fiber.StatusOK)
}
