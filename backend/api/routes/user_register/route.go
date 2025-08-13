package user_register

import (
	"trxd/api/auth"
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
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Username == "" || data.Email == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) < consts.MinPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.ShortPassword)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Username) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}
	if len(data.Email) > consts.MaxEmailLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongEmail)
	}

	if !consts.UserRegex.MatchString(data.Email) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidEmail)
	}

	user, err := RegisterUser(c.Context(), data.Username, data.Email, data.Password)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	sess, err := auth.Store.Get(c)
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
