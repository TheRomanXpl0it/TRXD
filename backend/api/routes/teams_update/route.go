package teams_update

import (
	"trxd/api/validator"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		// TODO: change name
		Country string `json:"country" validate:"team_country"`
		Image   string `json:"image" validate:"team_image"`
		Bio     string `json:"bio" validate:"team_bio"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Country == "" && data.Image == "" && data.Bio == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	tid := c.Locals("tid").(int32)
	err = UpdateTeam(c.Context(), tid, data.Country, data.Image, data.Bio)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
