package configs_get

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	configs, err := GetConfigs(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfigs, err)
	}

	return c.Status(fiber.StatusOK).JSON(configs)
}
