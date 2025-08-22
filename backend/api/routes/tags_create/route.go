package tags_create

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
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

	err := CreateTag(c.Context(), *data.ChallID, data.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return utils.Error(c, fiber.StatusConflict, consts.TagAlreadyExists)
			}
			if pqErr.Code == "23503" { // Foreign key violation error code
				return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
			}
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingTag, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
