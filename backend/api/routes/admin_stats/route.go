package admin_stats

import (
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	stats, err := GetAdminStats(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingStats)
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}
