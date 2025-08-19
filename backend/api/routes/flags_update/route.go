package flags_update

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
		Flag    string `json:"flag"`
		Regex   *bool  `json:"regex"`
		NewFlag string `json:"new_flag"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.Flag == "" || (data.Regex == nil && data.NewFlag == "") {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Flag) > consts.MaxFlagLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongFlag)
	}
	if len(data.NewFlag) > consts.MaxFlagLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongFlag)
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
