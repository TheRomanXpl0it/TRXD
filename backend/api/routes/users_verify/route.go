package users_verify

import (
	"encoding/json"
	"trxd/api/routes/users_register"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/jwt"

	"github.com/gofiber/fiber/v2"
)

func parseAndValidateToken(c *fiber.Ctx) (string, error) {
	token := c.Query("token", "")
	if token == "" {
		return "", utils.Error(c, fiber.StatusBadRequest, consts.InvalidToken)
	}

	claims, err := jwt.ParseAndValidateJWT(c.Context(), token)
	if err != nil {
		return "", utils.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return "", utils.Error(c, fiber.StatusUnauthorized, consts.InvalidToken)
	}

	return email, nil
}

func loginUser(c *fiber.Ctx, userID int32) (bool, error) {
	sess, err := db.Store.Get(c)
	if err != nil {
		return false, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", userID)

	err = sess.Save()
	if err != nil {
		return false, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return true, nil
}

func Route(c *fiber.Ctx) error {
	email, err := parseAndValidateToken(c)
	if err != nil || email == "" {
		return err
	}

	content, err := db.StorageGet(c.Context(), users_register.VerifyPrefix+email)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InternalServerError, err)
	}
	if content == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.InvalidToken)
	}

	var data users_register.Data
	err = json.Unmarshal([]byte(*content), &data)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InvalidJSON, err)
	}

	userID, err := users_register.RegisterUser(c, data)
	if err != nil || userID == -1 {
		return err
	}

	err = db.StorageDelete(c.Context(), users_register.VerifyPrefix+email)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InternalServerError, err)
	}

	success, err := loginUser(c, userID)
	if err != nil || !success {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
