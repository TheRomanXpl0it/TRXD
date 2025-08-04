package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func RegisterTeam(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) < consts.MinPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.ShortPassword)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Name) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}

	uid := c.Locals("uid").(int32)

	team, err := db.GetTeamFromUser(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if team != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyInTeam)
	}

	team, err = db.RegisterTeam(c.Context(), data.Name, data.Password, uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
	}
	if team == nil {
		return utils.Error(c, fiber.StatusConflict, consts.TeamAlreadyExists)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name": team.Name,
		"id":   team.ID,
	})
}

func JoinTeam(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) < consts.MinPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.ShortPassword)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Name) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}

	uid := c.Locals("uid").(int32)

	team, err := db.GetTeamFromUser(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if team != nil {
		return utils.Error(c, fiber.StatusConflict, consts.AlreadyInTeam)
	}

	team, err = db.JoinTeam(c.Context(), data.Name, data.Password, uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
	}
	if team == nil {
		return utils.Error(c, fiber.StatusConflict, consts.InvalidTeamCredentials)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name": team.Name,
		"id":   team.ID,
	})
}

func UpdateTeam(c *fiber.Ctx) error {
	var data struct {
		Nationality string `json:"nationality"`
		Image       string `json:"image"`
		Bio         string `json:"bio"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Nationality == "" && data.Image == "" && data.Bio == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if data.Nationality != "" && len(data.Nationality) > consts.MaxNationalityLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongNationality)
	}

	tid := c.Locals("tid").(int32)
	err := db.UpdateTeam(c.Context(), tid, data.Nationality, data.Image, data.Bio)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func ResetTeamPassword(c *fiber.Ctx) error {
	var data struct {
		TeamID *int32 `json:"team_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.TeamID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.TeamID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidTeamID)
	}

	newPassword, err := utils.GenerateRandPass()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorGeneratingPassword, err)
	}

	err = db.ResetTeamPassword(c.Context(), *data.TeamID, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingTeamPassword, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"new_password": newPassword,
	})
}

func GetTeams(c *fiber.Ctx) error {
	role := c.Locals("role")

	allData := false
	if role != nil {
		allData = utils.In(role.(db.UserRole), []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	}
	teamsData, err := db.GetTeams(c.Context(), allData)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(teamsData)
}

func GetTeam(c *fiber.Ctx) error {
	role := c.Locals("role")

	teamID, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidTeamID)
	}

	allData := false
	if !allData && role != nil {
		allData = utils.In(role.(db.UserRole), []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	}
	teamData, err := db.GetTeam(c.Context(), int32(teamID), allData, false)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}
	if teamData == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.TeamNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(teamData)
}
