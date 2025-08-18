package challenge_submit

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
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.Flag == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Flag) > consts.MaxFlagLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongFlag)
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	uid := c.Locals("uid").(int32)

	status, first_blood, err := SubmitFlag(c.Context(), uid, *data.ChallID, data.Flag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSubmittingFlag, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      status,
		"first_blood": first_blood,
	})
}
