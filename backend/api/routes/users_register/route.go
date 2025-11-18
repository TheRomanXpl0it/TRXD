package users_register

import (
	"fmt"
	"trxd/api/routes/teams_register"
	"trxd/api/validator"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/tde-nico/log"
)

func Route(c *fiber.Ctx) error {
	conf, err := db.GetConfig(c.Context(), "allow-register")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if conf != "true" {
		return utils.Error(c, fiber.StatusForbidden, consts.DisabledRegistrations)
	}

	uid := c.Locals("uid")
	if uid != nil {
		return utils.Error(c, fiber.StatusForbidden, consts.AlreadyRegistered)
	}

	var data struct {
		Name     string `json:"name" validate:"required,user_name"`
		Email    string `json:"email" validate:"required,user_email"`
		Password string `json:"password" validate:"required,password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
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

	user, err := RegisterUser(c.Context(), tx, data.Name, data.Email, data.Password)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	mode, err := db.GetConfig(c.Context(), "user-mode")
	if err != nil {
		log.Error("Failed to get user-mode config:", "err", err)
		mode = fmt.Sprint(consts.DefaultConfigs["user-mode"])
	}
	if mode == "true" {
		team, err := teams_register.RegisterTeam(c.Context(), tx, data.Name, data.Password, user.ID)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
		}
		if team == nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.TeamAlreadyExists, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCommittingTransaction, err)
	}

	sess, err := db.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", user.ID)

	err = sess.Save()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
