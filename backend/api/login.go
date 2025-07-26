package api

import (
	"trxd/db"

	"github.com/gofiber/fiber/v2"
)

func login(c *fiber.Ctx) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return apiError(c, fiber.StatusBadRequest, invalidJSON)
	}

	if data.Email == "" || data.Password == "" {
		return apiError(c, fiber.StatusBadRequest, missingRequiredFields)
	}
	if len(data.Password) > maxPasswordLength {
		return apiError(c, fiber.StatusBadRequest, longPassword)
	}
	if len(data.Email) > maxEmailLength {
		return apiError(c, fiber.StatusBadRequest, longEmail)
	}

	user, err := db.LoginUser(c.Context(), data.Email, data.Password)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorLoggingIn, err)
	}
	if user == nil {
		return apiError(c, fiber.StatusUnauthorized, invalidCredentials)
	}

	sess, err := store.Get(c)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorFetchingSession, err)
	}

	sess.Set("uid", user.ID)
	err = sess.Save()
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, errorSavingSession, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Name,
		"role":     user.Role,
	})
}
