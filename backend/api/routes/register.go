package routes

import (
	"trxd/api/auth"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	//TODO: check if configs["allow-register"] is true + make tests

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

	user, err := db.RegisterUser(c.Context(), data.Username, data.Email, data.Password)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Name,
		"role":     user.Role,
	})
}
