package users_register

import (
	"database/sql"
	"fmt"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/email"
	"trxd/utils/jwt"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

// TODO: move these into consts
const SUBJECT = "Email Verification Required"
const BODY_TEMPLATE = "Hello,\n\nTo Confirm your email address, please click the link below:\nhttp://%s/api/register?token=%s\n\nThank you!"

func verifyMailEnabled(c *fiber.Ctx) (bool, error) {
	server, err := db.GetConfig(c.Context(), "email-server")
	if err != nil {
		return false, err
	}
	if server == "" {
		return false, nil
	}

	return true, nil
}

func registerViaMail(c *fiber.Ctx, registerEmail string) error {
	valid, err := validator.Var(c, registerEmail, "required,user_email")
	if err != nil || !valid {
		return err
	}

	domain, err := db.GetConfig(c.Context(), "domain")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if domain == "" {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InvalidDomain)
	}

	exists, err := db.Sql.UserExistsByEmail(c.Context(), registerEmail)
	if err != nil {
		if err != sql.ErrNoRows {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
		}
	}
	if exists {
		return utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	err = email.InitEmailClientFromConfigs(c.Context())
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorInitializingEmailClient, err)
	}

	signed, err := jwt.GenerateJWT(c.Context(), jwt.Map{"email": registerEmail})
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSigningVerificationToken, err)
	}

	body := fmt.Sprintf(BODY_TEMPLATE, domain, signed)
	err = email.SendEmail(c.Context(), registerEmail, SUBJECT, body)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSendingVerificationEmail, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func parseAndValidateToken(c *fiber.Ctx, token string) (string, error) {
	valid, err := validator.Var(c, token, "required,jwt")
	if err != nil || !valid {
		return "", err
	}

	claims, err := jwt.ParseAndValidateJWT(c.Context(), token)
	if err != nil {
		return "", utils.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	emailAny, ok := claims["email"]
	if !ok || emailAny == nil {
		return "", utils.Error(c, fiber.StatusUnauthorized, consts.InvalidToken)
	}

	email, ok := emailAny.(string)
	if !ok || email == "" {
		return "", utils.Error(c, fiber.StatusUnauthorized, consts.InvalidToken)
	}

	valid, err = validator.Var(c, email, "required,user_email")
	if err != nil || !valid {
		return "", utils.Error(c, fiber.StatusUnauthorized, consts.InvalidEmail)
	}

	return email, nil
}
