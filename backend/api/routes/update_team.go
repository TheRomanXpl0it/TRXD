package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func UpdateTeam(c *fiber.Ctx) error {
	tid := c.Locals("tid")
	if tid != nil && tid.(int32) == -1 {
		return utils.Error(c, fiber.StatusForbidden, consts.Unauthorized)
	}

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

	err := db.UpdateTeam(c.Context(), tid.(int32), data.Nationality, data.Image, data.Bio)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
