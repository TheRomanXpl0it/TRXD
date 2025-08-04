package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tde-nico/log"
)

const PassLen = 10

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

func BytesToHex(data []byte) string {
	dataHex := make([]byte, len(data)*2)
	hex.Encode(dataHex, data)
	return string(dataHex)
}

func GenerateRandPass() (string, error) {
	data := make([]byte, PassLen)
	n, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	if n != PassLen {
		return "", fmt.Errorf("expected to read %d bytes, but got %d", PassLen, n)
	}
	return BytesToHex(data), nil
}

func Compare(a, b interface{}) error {
	expectedBytes, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal expected response: %v", err)
	}

	actualBytes, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal actual response: %v", err)
	}

	if string(expectedBytes) != string(actualBytes) {
		return fmt.Errorf("\nExpected:\n%s\nGot:\n%s", expectedBytes, actualBytes)
	}

	return nil
}
