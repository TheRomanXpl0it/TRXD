package teams_scoreboard_graph

import (
	"net/http"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	top, err := QueryTeamScoreboardGraph(c.Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, consts.ErrorFetchingScoreboardGraph, err)
	}

	return c.Status(fiber.StatusOK).JSON(top)
}
