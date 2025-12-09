package categories_update

import (
	"trxd/api/routes/categories_create"
	"trxd/api/routes/categories_delete"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name    string `json:"name" validate:"required,category_name"`
		NewName string `json:"new_name" validate:"required,category_name"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	cat, err := GetCategory(c.Context(), data.Name)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingCategory, err)
	}
	if cat == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.CategoryNotFound)
	}

	if data.NewName != "" && data.Name != data.NewName {
		category, err := categories_create.CreateCategory(c.Context(), data.NewName)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingCategory, err)
		}
		if category == nil {
			return utils.Error(c, fiber.StatusConflict, consts.CategoryAlreadyExists)
		}

		err = UpdateChallengesCategory(c.Context(), data.Name, data.NewName)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingCategory, err)
		}

		err = categories_delete.DeleteCategory(c.Context(), data.Name)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingCategory, err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
