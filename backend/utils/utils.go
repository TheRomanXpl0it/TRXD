package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tde-nico/log"
)

func In[T comparable](slice T, items []T) bool {
	for _, item := range items {
		if item == slice {
			return true
		}
	}
	return false
}

func Error(c *fiber.Ctx, status int, message string, err ...error) error {
	if len(err) != 0 {
		log.Error("API Error:", "desc", message, "err", err[0])
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}
