package crypto_utils

import (
	"bytes"
	"crypto/md5"
	"io"
	"trxd/utils"

	"golang.org/x/crypto/argon2"
)

const (
	TIME    = 1
	MEMORY  = 64 * 1024
	THREAD  = 4
	KEY_LEN = 32
)

func Hash(password string) (string, string, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return "", "", err
	}

	hashed := argon2.IDKey(
		[]byte(password),
		salt,
		TIME,
		MEMORY,
		THREAD,
		KEY_LEN,
	)

	hashedHex, err := utils.BytesToHex(hashed)
	if err != nil {
		return "", "", err
	}

	saltHex, err := utils.BytesToHex(salt)
	if err != nil {
		return "", "", err
	}

	return hashedHex, saltHex, nil
}

func Verify(password string, saltHex string, hashHex string) (bool, error) {
	hash, err := utils.HextoBytes(hashHex)
	if err != nil {
		return false, err
	}

	salt, err := utils.HextoBytes(saltHex)
	if err != nil {
		return false, err
	}

	hashed := argon2.IDKey(
		[]byte(password),
		salt,
		TIME,
		MEMORY,
		THREAD,
		KEY_LEN,
	)

	return bytes.Equal(hashed, hash), nil
}

func HashFile(file io.Reader) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	h := md5.New()
	n, err := h.Write(data)
	if err != nil {
		return "", err
	}
	if n != len(data) {
		return "", io.ErrShortWrite
	}

	hash := h.Sum(nil)
	hashHex, err := utils.BytesToHex(hash)
	if err != nil {
		return "", err
	}

	return hashHex, nil
}
