package users_register

import (
	"trxd/api/middlewares"
	"trxd/api/validator"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	conf, err := db.GetConfig(c.Context(), "allow-register")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if conf == nil || conf.Value != "true" {
		return utils.Error(c, fiber.StatusForbidden, consts.DisabledRegistration)
	}

	uid := c.Locals("uid")
	if uid != nil {
		return utils.Error(c, fiber.StatusForbidden, consts.AlreadyRegistered)
	}

	var data struct {
		Name     string `json:"name" validate:"required,user_name"`
		Email    string `json:"email" validate:"required,user_email"`
		Password string `json:"password" validate:"required,user_password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	if !consts.UserRegex.MatchString(data.Email) { // TODO: put into the validator
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidEmail)
	}

	user, err := RegisterUser(c.Context(), data.Name, data.Email, data.Password)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	sess, err := middlewares.Store.Get(c)
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
