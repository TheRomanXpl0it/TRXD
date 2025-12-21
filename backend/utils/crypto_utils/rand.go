package crypto_utils

import (
	"crypto/rand"
	"fmt"
	"trxd/utils"
)

const PassLen = 10
const SaltLen = 16

func GeneratePassword() (string, error) {
	data := make([]byte, PassLen)
	n, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	if n != PassLen {
		return "", fmt.Errorf("expected to read %d bytes, but got %d", PassLen, n)
	}

	password, err := utils.BytesToHex(data)
	if err != nil {
		return "", err
	}

	return password, nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	n, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	if n != SaltLen {
		return nil, fmt.Errorf("expected to read %d bytes, but got %d", SaltLen, n)
	}

	return salt, nil
}
