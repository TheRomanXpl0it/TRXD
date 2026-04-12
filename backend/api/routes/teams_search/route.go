package teams_search

import (
	"trxd/api/routes/teams_get"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func FetchEmail(c *fiber.Ctx, role interface{}, email string) (*teams_get.TeamData, error) {
	if role == nil || role.(sqlc.UserRole) != sqlc.UserRoleAdmin {
		return nil, utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	valid, err := validator.Var(c, email, "user_email")
	if err != nil || !valid {
		return nil, err
	}

	teamData, err := GetTeamByEmail(c.Context(), email)
	if err != nil {
		return nil, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if teamData == nil {
		return nil, utils.Error(c, fiber.StatusNotFound, consts.TeamNotFound)
	}

	return teamData, nil
}

func FetchName(c *fiber.Ctx, tid interface{}, role interface{}, name string) (*teams_get.TeamData, error) {
	valid, err := validator.Var(c, name, "team_name")
	if err != nil || !valid {
		return nil, err
	}

	allData := false
	if role != nil {
		allData = utils.In(role.(sqlc.UserRole), []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	}
	teamData, err := GetTeamByName(c.Context(), name, tid, allData)
	if err != nil {
		return nil, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if teamData == nil {
		return nil, utils.Error(c, fiber.StatusNotFound, consts.TeamNotFound)
	}

	return teamData, nil
}

func Route(c *fiber.Ctx) error {
	tid := c.Locals("tid")
	role := c.Locals("role")

	teamName := c.Query("name")
	teamEmail := c.Query("email")

	if teamEmail == "" && teamName == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}

	var teamData *teams_get.TeamData
	var err error

	if teamEmail != "" {
		teamData, err = FetchEmail(c, role, teamEmail)
		if err != nil || teamData == nil {
			return err
		}
	} else if teamName != "" {
		teamData, err = FetchName(c, tid, role, teamName)
		if err != nil || teamData == nil {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(teamData)
}
