package flags_delete

import (
	"trxd/api/validator"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id" validate:"required,challenge_id"`
		Flag    string `json:"flag" validate:"required,flag_flag"`
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

	err = DeleteFlag(c.Context(), *data.ChallID, data.Flag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingFlag, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
