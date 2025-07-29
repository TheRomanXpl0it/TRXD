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
		return apiError(c, fiber.StatusBadRequest, InvalidJSON)
	}

	if data.Email == "" || data.Password == "" {
		return apiError(c, fiber.StatusBadRequest, MissingRequiredFields)
	}
	if len(data.Password) > MaxPasswordLength {
		return apiError(c, fiber.StatusBadRequest, LongPassword)
	}
	if len(data.Email) > MaxEmailLength {
		return apiError(c, fiber.StatusBadRequest, LongEmail)
	}

	user, err := db.LoginUser(c.Context(), data.Email, data.Password)
	if err != nil {
		return apiError(c, fiber.StatusInternalServerError, ErrorLoggingIn, err)
	}
	if user == nil {
		return apiError(c, fiber.StatusUnauthorized, InvalidCredentials)
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
