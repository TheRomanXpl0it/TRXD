package tags_update

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
		OldName string `json:"old_name"`
		NewName string `json:"new_name"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.OldName == "" || data.NewName == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.ChallID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallengeID)
	}
	if len(data.OldName) > consts.MaxTagNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongTagName)
	}
	if len(data.NewName) > consts.MaxTagNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongTagName)
	}

	err := UpdateTag(c.Context(), *data.ChallID, data.OldName, data.NewName)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingTag, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
