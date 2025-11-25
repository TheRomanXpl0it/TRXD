package categories_create

import (
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name string `json:"name" validate:"required,category_name"`
		Icon string `json:"icon" validate:"required,category_icon"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	category, err := CreateCategory(c.Context(), data.Name, data.Icon)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingCategory, err)
	}
	if category == nil {
		return utils.Error(c, fiber.StatusConflict, consts.CategoryAlreadyExists)
	}

	return c.SendStatus(fiber.StatusOK)
}
