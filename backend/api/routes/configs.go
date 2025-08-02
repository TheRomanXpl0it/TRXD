package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func UpdateConfig(c *fiber.Ctx) error {
	var data struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Key == "" || data.Value == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}

	conf, err := db.GetConfig(c.Context(), data.Key)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingConfig, err)
	}
	if conf == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ConfigNotFound)
	}

	if conf.Value == data.Value {
		return c.SendStatus(fiber.StatusOK) // No change needed
	}

	err = db.UpdateConfig(c.Context(), data.Key, data.Value)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingConfig, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
