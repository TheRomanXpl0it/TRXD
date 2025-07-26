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
		return apiError(c, fiber.StatusBadRequest, invalidJSON)
	}

	if data.Username == "" || data.Email == "" || data.Password == "" {
		return apiError(c, fiber.StatusBadRequest, missingRequiredFields)
	}
	if len(data.Password) < minPasswordLength {
		return apiError(c, fiber.StatusBadRequest, shortPassword)
	}
	if len(data.Password) > maxPasswordLength {
		return apiError(c, fiber.StatusBadRequest, longPassword)
	}
	if len(data.Username) > maxNameLength {
		return apiError(c, fiber.StatusBadRequest, longName)
	}
	if len(data.Email) > maxEmailLength {
		return apiError(c, fiber.StatusBadRequest, longEmail)
	}

	if !UserRegex.MatchString(data.Email) {
		return apiError(c, fiber.StatusBadRequest, invalidEmail)
	}

	user, err := db.RegisterUser(c.Context(), data.Username, data.Email, data.Password)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorRegisteringUser, err)
	}
	if user == nil {
		return apiError(c, fiber.StatusConflict, userAlreadyExists)
	}

	sess, err := store.Get(c)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorFetchingSession, err)
	}

	sess.Set("uid", user.ID)
	err = sess.Save()
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorSavingSession)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Name,
		"role":     user.Role,
	})
}
