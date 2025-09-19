package teams_get

import (
	"trxd/api/validator"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role")

	teamID, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidTeamID)
	}
	valid, err := validator.Var(c, teamID, "team_id")
	if err != nil || !valid {
		return err
	}

	allData := false
	if !allData && role != nil {
		allData = utils.In(role.(sqlc.UserRole), []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	}
	teamData, err := GetTeam(c.Context(), int32(teamID), allData)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if teamData == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.TeamNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(teamData)
}
