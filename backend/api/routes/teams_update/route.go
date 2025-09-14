package teams_update

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Country string `json:"country"`
		Image   string `json:"image"`
		Bio     string `json:"bio"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Country == "" && data.Image == "" && data.Bio == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if data.Country != "" && len(data.Country) > consts.MaxCountryLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongCountry)
	}
	if data.Image != "" && len(data.Image) > consts.MaxImageLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongImage)
	}
	if data.Bio != "" && len(data.Bio) > consts.MaxBioLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongBio)
	}

	tid := c.Locals("tid").(int32)
	err := UpdateTeam(c.Context(), tid, data.Country, data.Image, data.Bio)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
