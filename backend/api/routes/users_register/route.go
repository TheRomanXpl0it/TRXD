package users_register

import (
	"trxd/api/routes/teams_register"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

type Data struct {
	Name     string `json:"name" validate:"required,user_name"`
	Email    string `json:"email" validate:"required,user_email"`
	Password string `json:"password" validate:"required,password"`
	JWT      string `json:"token" validate:"omitempty,jwt"`
}

func CanRegister(c *fiber.Ctx) (bool, error) {
	conf, err := db.GetConfig(c.Context(), "allow-register")
	if err != nil {
		return false, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if conf != "true" {
		return false, utils.Error(c, fiber.StatusForbidden, consts.DisabledRegistrations)
	}

	uid := c.Locals("uid")
	if uid != nil {
		return false, utils.Error(c, fiber.StatusForbidden, consts.AlreadyRegistered)
	}

	return true, nil
}

func RegisterUser(c *fiber.Ctx, data Data) (int32, error) {
	tx, err := db.BeginTx(c.Context())
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorBeginningTransaction, err)
	}
	defer db.Rollback(tx)

	user, err := DBRegisterUser(c.Context(), tx, data.Name, data.Email, data.Password)
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if user == nil {
		return -1, utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	mode, err := db.GetConfig(c.Context(), "user-mode")
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if mode == "true" {
		team, err := teams_register.RegisterTeam(c.Context(), tx, data.Name, data.Password, user.ID)
		if err != nil {
			return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
		}
		if team == nil {
			return -1, utils.Error(c, fiber.StatusInternalServerError, consts.TeamAlreadyExists, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCommittingTransaction, err)
	}

	return user.ID, nil
}

func loginUser(c *fiber.Ctx, userID int32) (bool, error) {
	sess, err := db.Store.Get(c)
	if err != nil {
		return false, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", userID)

	err = sess.Save()
	if err != nil {
		return false, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return true, nil
}

func Route(c *fiber.Ctx) error {
	canRegister, err := CanRegister(c)
	if err != nil || !canRegister {
		return err
	}

	var data Data
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	enabled, err := verifyMailEnabled(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}

	if enabled {
		if data.JWT == "" {
			return registerViaMail(c, data.Email)
		}

		mail, err := parseAndValidateToken(c, data.JWT)
		if err != nil || mail == "" {
			return err
		}

		data.Email = mail
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	userID, err := RegisterUser(c, data)
	if err != nil || userID == -1 {
		return err
	}

	success, err := loginUser(c, userID)
	if err != nil || !success {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
