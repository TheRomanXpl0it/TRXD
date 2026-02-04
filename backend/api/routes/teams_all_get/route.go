package teams_all_get

import (
	"math"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	offset := c.QueryInt("offset", 0)
	if offset < 0 || offset > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	limit := c.QueryInt("limit", 0)
	if limit < 0 || limit > math.MaxInt32 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidParam)
	}

	totalTeams, teamsData, err := GetTeams(c.Context(), int32(offset), int32(limit))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total": totalTeams,
		"teams": teamsData,
	})
}
