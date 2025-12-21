package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

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
	if c == nil {
		return errors.New(message)
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}

func BytesToHex(data []byte) (string, error) {
	dataHex := make([]byte, len(data)*2)

	n := hex.Encode(dataHex, data)
	if n != len(data)*2 {
		return "", errors.New("failed to encode bytes to hex")
	}

	return string(dataHex), nil
}

func HextoBytes(dataHex string) ([]byte, error) {
	data := make([]byte, (len(dataHex)+1)/2)

	_, err := hex.Decode(data, []byte(dataHex))
	if err != nil {
		return nil, err
	}

	return data, nil
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
