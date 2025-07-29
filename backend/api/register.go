package api

import (
	"trxd/db"

	"github.com/gofiber/fiber/v2"
)

func register(c *fiber.Ctx) error {
	var data struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return apiError(c, fiber.StatusBadRequest, InvalidJSON)
	}

	if data.Username == "" || data.Email == "" || data.Password == "" {
		return apiError(c, fiber.StatusBadRequest, MissingRequiredFields)
	}
	if len(data.Password) < MinPasswordLength {
		return apiError(c, fiber.StatusBadRequest, ShortPassword)
	}
	if len(data.Password) > MaxPasswordLength {
		return apiError(c, fiber.StatusBadRequest, LongPassword)
	}
	if len(data.Username) > MaxNameLength {
		return apiError(c, fiber.StatusBadRequest, LongName)
	}
	if len(data.Email) > MaxEmailLength {
		return apiError(c, fiber.StatusBadRequest, LongEmail)
	}

	if !UserRegex.MatchString(data.Email) {
		return apiError(c, fiber.StatusBadRequest, InvalidEmail)
	}

	user, err := db.RegisterUser(c.Context(), data.Username, data.Email, data.Password)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorRegisteringUser, err)
	}
	if user == nil {
		return apiError(c, fiber.StatusConflict, UserAlreadyExists)
	}

	sess, err := store.Get(c)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorFetchingSession, err)
	}

	sess.Set("uid", user.ID)
	err = sess.Save()
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorSavingSession, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Name,
		"role":     user.Role,
	})
}
