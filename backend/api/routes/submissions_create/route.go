package submissions_create

import (
	"strings"
	"trxd/api/validator"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/discord"

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

	uid := c.Locals("uid").(int32)
	data.Flag = strings.TrimSpace(data.Flag)

	status, first_blood, err := SubmitFlag(c.Context(), uid, *data.ChallID, data.Flag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSubmittingFlag, err)
	}

	if first_blood {
		discord.BroadcastFirstBlood(c.Context(), challenge, uid)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      status,
		"first_blood": first_blood,
	})
}
