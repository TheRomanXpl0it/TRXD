package users_update

import (
	"fmt"
	"trxd/api/routes/teams_update"
	"trxd/api/validator"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/tde-nico/log"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name    string  `json:"name" validate:"user_name"`
		Country *string `json:"country" validate:"omitempty,country"`
		Image   *string `json:"image" validate:"omitempty,image_url"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" && data.Country == nil && data.Image == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	uid := c.Locals("uid").(int32)

	tx, err := db.BeginTx(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorBeginningTransaction, err)
	}
	defer tx.Rollback()

	err = UpdateUser(c.Context(), tx, uid, data.Name, data.Country, data.Image)
	if err != nil {
		if err.Error() == "[name already taken]" {
			return utils.Error(c, fiber.StatusConflict, consts.NameAlreadyTaken)
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	mode, err := db.GetConfig(c.Context(), "user-mode")
	if err != nil {
		log.Error("Failed to get user-mode config:", "err", err)
		mode = fmt.Sprint(consts.DefaultConfigs["user-mode"])
	}
	if mode == "true" {
		tid := c.Locals("tid").(int32)

		err = teams_update.UpdateTeam(c.Context(), tx, tid, data.Name, data.Country, data.Image, nil)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingTeam, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCommittingTransaction, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
