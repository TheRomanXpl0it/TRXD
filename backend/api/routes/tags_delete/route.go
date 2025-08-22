package tags_delete

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
		Name    string `json:"name"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.Name == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Name) > consts.MaxTagNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongTagName)
	}

	err := DeleteTag(c.Context(), *data.ChallID, data.Name)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingTag, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
