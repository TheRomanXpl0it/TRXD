package categories_delete

import (
	"trxd/api/validator"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Category string `json:"category" validate:"required,category_name"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	err = DeleteCategory(c.Context(), data.Category)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingCategory, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
