package api

import (
	"trxd/db"

	"github.com/gofiber/fiber/v2"
)

func joinTeam(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return apiError(c, fiber.StatusBadRequest, invalidJSON)
	}

	if data.Name == "" || data.Password == "" {
		return apiError(c, fiber.StatusBadRequest, missingRequiredFields)
	}
	if len(data.Password) < minPasswordLength {
		return apiError(c, fiber.StatusBadRequest, shortPassword)
	}
	if len(data.Password) > maxPasswordLength {
		return apiError(c, fiber.StatusBadRequest, longPassword)
	}
	if len(data.Name) > maxNameLength {
		return apiError(c, fiber.StatusBadRequest, longName)
	}

	uid := c.Locals("uid")

	team, err := db.GetTeamFromUser(c.Context(), uid.(int32))
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorFetchingTeam, err)
	}
	if team != nil {
		return apiError(c, fiber.StatusConflict, alreadyInTeam)
	}

	team, err = db.JoinTeam(c.Context(), data.Name, data.Password, uid.(int32))
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorRegisteringTeam, err)
	}
	if team == nil {
		return apiError(c, fiber.StatusConflict, invalidTeamCredentials)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name": team.Name,
	})
}
