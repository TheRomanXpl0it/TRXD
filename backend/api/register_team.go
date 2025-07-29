package api

import (
	"trxd/db"

	"github.com/gofiber/fiber/v2"
)

func registerTeam(c *fiber.Ctx) error {
	var data struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return apiError(c, fiber.StatusBadRequest, InvalidJSON)
	}

	if data.Name == "" || data.Password == "" {
		return apiError(c, fiber.StatusBadRequest, MissingRequiredFields)
	}
	if len(data.Password) < MinPasswordLength {
		return apiError(c, fiber.StatusBadRequest, ShortPassword)
	}
	if len(data.Password) > MaxPasswordLength {
		return apiError(c, fiber.StatusBadRequest, LongPassword)
	}
	if len(data.Name) > MaxNameLength {
		return apiError(c, fiber.StatusBadRequest, LongName)
	}

	uid := c.Locals("uid")

	team, err := db.GetTeamFromUser(c.Context(), uid.(int32))
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorFetchingTeam, err)
	}
	if team != nil {
		return apiError(c, fiber.StatusConflict, AlreadyInTeam)
	}

	team, err = db.RegisterTeam(c.Context(), data.Name, data.Password, uid.(int32))
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorRegisteringTeam, err)
	}
	if team == nil {
		return apiError(c, fiber.StatusConflict, TeamAlreadyExists)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name": team.Name,
	})
}
