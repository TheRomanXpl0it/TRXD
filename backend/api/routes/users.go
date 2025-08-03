package routes

import (
	"trxd/api/auth"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	conf, err := db.GetConfig(c.Context(), "allow-register")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	if conf == nil || conf.Value != "true" {
		return utils.Error(c, fiber.StatusForbidden, consts.DisabledRegistration)
	}

	var data struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Username == "" || data.Email == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) < consts.MinPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.ShortPassword)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Username) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}
	if len(data.Email) > consts.MaxEmailLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongEmail)
	}

	if !consts.UserRegex.MatchString(data.Email) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidEmail)
	}

	user, err := db.RegisterUser(c.Context(), data.Username, data.Email, data.Password)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorRegisteringUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusConflict, consts.UserAlreadyExists)
	}

	sess, err := auth.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", user.ID)
	err = sess.Save()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Name,
		"role":     user.Role,
	})
}

func Login(c *fiber.Ctx) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Email == "" || data.Password == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Password) > consts.MaxPasswordLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongPassword)
	}
	if len(data.Email) > consts.MaxEmailLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongEmail)
	}

	user, err := db.LoginUser(c.Context(), data.Email, data.Password)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorLoggingIn, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.InvalidCredentials)
	}

	sess, err := auth.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	sess.Set("uid", user.ID)
	err = sess.Save()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingSession, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Name,
		"role":     user.Role,
	})
}

func Logout(c *fiber.Ctx) error {
	sess, err := auth.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	err = sess.Destroy()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDestroyingSession, err)
	}

	return c.Status(fiber.StatusOK).SendString("")
}

func UpdateUser(c *fiber.Ctx) error {
	var data struct {
		Name        string `json:"name"`
		Nationality string `json:"nationality"`
		Image       string `json:"image"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" && data.Nationality == "" && data.Image == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if data.Name != "" && len(data.Name) > consts.MaxNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongName)
	}
	if data.Nationality != "" && len(data.Nationality) > consts.MaxNationalityLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongNationality)
	}

	uid := c.Locals("uid").(int32)

	err := db.UpdateUser(c.Context(), uid, data.Name, data.Nationality, data.Image)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingUser, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func ResetUserPassword(c *fiber.Ctx) error {
	var data struct {
		UserID *int32 `json:"user_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.UserID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.UserID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidUserID)
	}

	newPassword, err := utils.GenerateRandPass()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorGeneratingPassword, err)
	}

	err = db.ResetUserPassword(c.Context(), *data.UserID, newPassword)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorResettingUserPassword, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"new_password": newPassword})
}

func Info(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)

	user, err := db.GetUserByID(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser)
	}

	team, err := db.GetTeamFromUser(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingTeam, err)
	}

	teamID := int32(-1)
	teamName := ""
	if team != nil {
		teamID = team.ID
		teamName = team.Name
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username":  user.Name,
		"role":      user.Role,
		"team_id":   teamID,
		"team_name": teamName,
	})
}

// TODO: tests
func GetUsers(c *fiber.Ctx) error {
	role := c.Locals("role").(db.UserRole)

	allData := utils.In(role, []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	usersData, err := db.GetUsers(c.Context(), allData)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	return c.Status(fiber.StatusOK).JSON(usersData)
}

// TODO: tests
func GetUser(c *fiber.Ctx) error {
	role := c.Locals("role").(db.UserRole)

	userID, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidUserID)
	}

	allData := utils.In(role, []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	userData, err := db.GetUser(c.Context(), int32(userID), allData, false)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if userData == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.UserNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(userData)
}
