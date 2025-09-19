package flags_update

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
		Regex   *bool  `json:"regex"`
		NewFlag string `json:"new_flag" validate:"flag_flag"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Regex == nil && data.NewFlag == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
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

	ok, err := UpdateFlag(c.Context(), *data.ChallID, data.Flag, data.Regex, data.NewFlag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingFlag, err)
	}
	if !ok {
		return utils.Error(c, fiber.StatusConflict, consts.FlagAlreadyExists)
	}

	return c.SendStatus(fiber.StatusOK)
}
