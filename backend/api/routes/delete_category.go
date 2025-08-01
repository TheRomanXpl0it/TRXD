package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func DeleteCategory(c *fiber.Ctx) error {
	var data struct {
		Category string `json:"category"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Category == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Category) > consts.MaxCategoryLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongCategory)
	}

	err := db.DeleteCategory(c.Context(), data.Category)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingCategory, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
