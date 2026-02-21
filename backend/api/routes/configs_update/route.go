package configs_update

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"trxd/utils/log"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Key   string  `json:"key" validate:"required"`
		Value *string `json:"value" validate:"required"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	conf, err := db.GetCompleteConfig(c.Context(), data.Key)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingConfig, err)
	}
	if conf == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ConfigNotFound)
	}

	if conf.Value == *data.Value { // No change needed
		return c.SendStatus(fiber.StatusOK)
	}

	err = db.UpdateConfig(c.Context(), data.Key, *data.Value)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingConfig, err)
	}

	if data.Key == "user-mode" {
		log.Warn("Shutting down server to apply user-mode change")
		go func(app *fiber.App) {
			err := app.Shutdown()
			if err != nil {
				log.Error("Error shutting down server", "err", err)
			}
		}(c.App())
	}

	return c.SendStatus(fiber.StatusOK)
}
