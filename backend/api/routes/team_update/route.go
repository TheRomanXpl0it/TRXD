package team_update

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Nationality string `json:"nationality"`
		Image       string `json:"image"`
		Bio         string `json:"bio"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Nationality == "" && data.Image == "" && data.Bio == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if data.Nationality != "" && len(data.Nationality) > consts.MaxNationalityLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongNationality)
	}

	tid := c.Locals("tid").(int32)
	err := UpdateTeam(c.Context(), tid, data.Nationality, data.Image, data.Bio)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
