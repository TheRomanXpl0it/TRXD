package teams_update

import (
	"trxd/api/validator"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name    string  `json:"name" validate:"team_name"`
		Country *string `json:"country" validate:"omitempty,country"`
		Image   *string `json:"image" validate:"omitempty,image_url"`
		Bio     *string `json:"bio" validate:"omitempty,team_bio"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" && data.Country == nil && data.Image == nil && data.Bio == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	tx, err := db.BeginTx(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorBeginningTransaction, err)
	}
	defer tx.Rollback()

	tid := c.Locals("tid").(int32)
	err = UpdateTeam(c.Context(), tx, tid, data.Name, data.Country, data.Image, data.Bio)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return utils.Error(c, fiber.StatusConflict, consts.NameAlreadyTaken)
			}
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingTeam, err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCommittingTransaction, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
