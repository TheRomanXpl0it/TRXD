package users_update

import (
	"trxd/api/validator"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name    string `json:"name" validate:"user_name"`
		Country string `json:"country" validate:"user_country"`
		Image   string `json:"image" validate:"user_image"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" && data.Country == "" && data.Image == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	uid := c.Locals("uid").(int32)

	err = UpdateUser(c.Context(), uid, data.Name, data.Country, data.Image)
	if err != nil {
		if err.Error() == "[name already taken]" {
			return utils.Error(c, fiber.StatusConflict, consts.NameAlreadyTaken)
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
