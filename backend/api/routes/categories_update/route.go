package categories_update

import (
	"trxd/api/routes/categories_create"
	"trxd/api/routes/categories_delete"
	"trxd/api/validator"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name    string `json:"name" validate:"required,category_name"`
		NewName string `json:"new_name" validate:"category_name"`
		NewIcon string `json:"new_icon" validate:"category_name"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.NewName == "" && data.NewIcon == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
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
		icon := cat.Icon
		if icon != data.NewIcon {
			icon = data.NewIcon
		}

		category, err := categories_create.CreateCategory(c.Context(), data.NewName, icon)
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
	} else {
		err := UpdateCategoryIcon(c.Context(), data.Name, data.NewIcon)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingCategory, err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
