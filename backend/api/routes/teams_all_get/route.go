package teams_all_get

import (
	"math"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	if start < 0 || start > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	end := c.QueryInt("end", 0)
	if end < 0 || end > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	if end != 0 && end <= start {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	teamsData, err := GetTeams(c.Context(), int32(start), int32(end))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(teamsData)
}
