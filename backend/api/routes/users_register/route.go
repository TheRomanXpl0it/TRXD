package users_register

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"trxd/api/routes/teams_register"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/email"
	"trxd/utils/jwt"
	"trxd/validator"

	"trxd/utils/log"

	"github.com/gofiber/fiber/v2"
)

type Data struct {
	Name     string `json:"name" validate:"required,user_name"`
	Email    string `json:"email" validate:"required,user_email"`
	Password string `json:"password" validate:"required,password"`
}

// TODO: move these into consts
const VerifyPrefix = "verify-"
const SUBJECT = "Email Verification Required"
const BODY_TEMPLATE = "Hello %s,\n\nTo Confirm your email address, please click the link below:\nhttp://%s/api/verify?token=%s\n\nThank you!"

func SetNXUserData(ctx context.Context, data Data) (bool, error) {
	content, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	res, err := db.StorageSetNX(ctx, VerifyPrefix+data.Email, string(content), 1*time.Minute)
	if err != nil {
		return false, err
	}
	if !res {
		return false, nil
	}

	return true, nil
}

func registerViaMail(c *fiber.Ctx, data Data) (bool, error) {
	server, err := db.GetConfig(c.Context(), "email-server")
	if err != nil {
		return false, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if server == "" {
		return false, nil
	}

	domain, err := db.GetConfig(c.Context(), "domain")
	if err != nil {
		return true, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if domain == "" {
		return true, utils.Error(c, fiber.StatusInternalServerError, consts.InvalidDomain)
	}

	exists, err := db.Sql.UserExistsByEmail(c.Context(), data.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			return true, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
		}
	}
	if exists {
		return true, utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	err = email.InitEmailClientFromConfigs(c.Context())
	if err != nil {
		return true, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorInitializingEmailClient, err)
	}

	success, err := SetNXUserData(c.Context(), data)
	if err != nil {
		return true, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if !success {
		// feature: make retries
		return true, utils.Error(c, fiber.StatusTooManyRequests, consts.VerificationAlreadySent)
	}

	signed, err := jwt.GenerateJWT(c.Context(), jwt.Map{"email": data.Email})
	if err != nil {
		return true, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSigningVerificationToken, err)
	}

	body := fmt.Sprintf(BODY_TEMPLATE, data.Name, domain, signed)
	err = email.SendEmail(c.Context(), data.Email, SUBJECT, body)
	if err != nil {
		return true, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSendingVerificationEmail, err)
	}

	return true, c.SendStatus(fiber.StatusOK)
}

func RegisterUser(c *fiber.Ctx, data Data) (int32, error) {
	tx, err := db.BeginTx(c.Context())
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorBeginningTransaction, err)
	}
	defer db.Rollback(tx)

	user, err := DBRegisterUser(c.Context(), tx, data.Name, data.Email, data.Password)
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if user == nil {
		return -1, utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	mode, err := db.GetConfig(c.Context(), "user-mode")
	if err != nil {
		log.Error("Failed to get user-mode config:", "err", err)
		mode = fmt.Sprint(consts.DefaultConfigs["user-mode"])
	}
	if mode == "true" {
		team, err := teams_register.RegisterTeam(c.Context(), tx, data.Name, data.Password, user.ID)
		if err != nil {
			return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringTeam, err)
		}
		if team == nil {
			return -1, utils.Error(c, fiber.StatusInternalServerError, consts.TeamAlreadyExists, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCommittingTransaction, err)
	}

	return user.ID, nil
}

func Route(c *fiber.Ctx) error {
	conf, err := db.GetConfig(c.Context(), "allow-register")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if conf != "true" {
		return utils.Error(c, fiber.StatusForbidden, consts.DisabledRegistrations)
	}

	uid := c.Locals("uid")
	if uid != nil {
		return utils.Error(c, fiber.StatusForbidden, consts.AlreadyRegistered)
	}

	var data Data
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	// TODO: hash password here

	enabled, err := registerViaMail(c, data)
	if err != nil || enabled {
		return err
	}

	userID, err := RegisterUser(c, data)
	if err != nil || userID == -1 {
		return err
	}

	sess, err := db.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", userID)

	err = sess.Save()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
