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

	hashedHex := utils.BytesToHex(hashed)
	saltHex := utils.BytesToHex(salt)

	return hashedHex, saltHex, nil
}

func Verify(password string, saltHex string, hashHex string) bool {
	hash := utils.HextoBytes(hashHex)
	salt := utils.HextoBytes(saltHex)

	hashed := argon2.IDKey(
		[]byte(password),
		salt,
		TIME,
		MEMORY,
		THREAD,
		KEY_LEN,
	)

	return bytes.Equal(hashed, hash)
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
	return utils.BytesToHex(hash), nil
}
