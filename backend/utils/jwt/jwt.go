package jwt

import (
	"context"
	"errors"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/golang-jwt/jwt/v5"
)

type Map = jwt.MapClaims

var signingMethod = jwt.SigningMethodHS256

func fetchKey(ctx context.Context) ([]byte, error) {
	secretHex, err := db.GetConfig(ctx, "jwt-secret")
	if err != nil {
		return nil, err
	}
	if secretHex == "" {
		return nil, errors.New(consts.InvalidJWTSecret)
	}

	secret, err := utils.HextoBytes(secretHex)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func GenerateJWT(ctx context.Context, data Map) (string, error) {
	key, err := fetchKey(ctx)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(signingMethod, data)
	signed, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ParseAndValidateJWT(ctx context.Context, tokenString string) (Map, error) {
	key, err := fetchKey(ctx)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(consts.InvalidSigningMethod)
		}

		if token.Method != signingMethod {
			return nil, errors.New(consts.InvalidSigningAlgorithm)
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New(consts.InvalidToken)
	}

	return claims, nil
}
