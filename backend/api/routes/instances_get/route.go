package instances_get

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	instances, err := GetInstances(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingInstances, err)
	}

	return c.Status(fiber.StatusOK).JSON(instances)
}
