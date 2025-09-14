package users_update

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name    string `json:"name"`
		Country string `json:"country"`
		Image   string `json:"image"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" && data.Country == "" && data.Image == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if data.Name != "" && len(data.Name) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}
	if data.Country != "" && len(data.Country) > consts.MaxCountryLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongCountry)
	}
	if data.Image != "" && len(data.Image) > consts.MaxImageLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongImage)
	}

	uid := c.Locals("uid").(int32)

	err := UpdateUser(c.Context(), uid, data.Name, data.Country, data.Image)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
