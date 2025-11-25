package submissions_create

import (
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/discord"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id" validate:"required,id"`
		Flag    string `json:"flag" validate:"required,flag"`
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
	role := c.Locals("role").(sqlc.UserRole)
	if role == sqlc.UserRolePlayer && challenge.Hidden {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	data.Flag = strings.TrimSpace(data.Flag)

	status, first_blood, err := SubmitFlag(c.Context(), uid, role, *data.ChallID, data.Flag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSubmittingFlag, err)
	}

	if first_blood {
		go discord.BroadcastFirstBlood(c.Context(), challenge, uid)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      status,
		"first_blood": first_blood,
	})
}
