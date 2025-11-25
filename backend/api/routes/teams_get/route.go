package teams_get

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	role := c.Locals("role")

	teamIDInt, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidTeamID)
	}
	teamID := int32(teamIDInt)
	valid, err := validator.Var(c, teamID, "id")
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
