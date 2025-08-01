package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func CreateCategory(c *fiber.Ctx) error {
	var data struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" || data.Icon == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Name) > consts.MaxCategoryLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongCategory)
	}
	if len(data.Icon) > consts.MaxIconLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongIcon)
	}

	category, err := db.CreateCategory(c.Context(), data.Name, data.Icon)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingCategory, err)
	}
	if category == nil {
		return utils.Error(c, fiber.StatusConflict, consts.CategoryAlreadyExists)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name": category.Name,
		"icon": category.Icon,
	})
}
