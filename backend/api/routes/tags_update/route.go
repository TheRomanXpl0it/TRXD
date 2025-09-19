package tags_update

import (
	"trxd/api/validator"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id" validate:"required,challenge_id"`
		OldName string `json:"old_name" validate:"required,tag_name"`
		NewName string `json:"new_name" validate:"required,tag_name"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	err = UpdateTag(c.Context(), *data.ChallID, data.OldName, data.NewName)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingTag, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
