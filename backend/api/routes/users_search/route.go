package users_search

import (
	"trxd/api/routes/users_get"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func FetchEmail(c *fiber.Ctx, role interface{}, email string) (*users_get.UserData, error) {
	if role == nil || role.(sqlc.UserRole) != sqlc.UserRoleAdmin {
		return nil, utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	valid, err := validator.Var(c, email, "user_email")
	if err != nil || !valid {
		return nil, err
	}

	userData, err := GetUserByEmail(c.Context(), email)
	if err != nil {
		return nil, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if userData == nil {
		return nil, utils.Error(c, fiber.StatusNotFound, consts.UserNotFound)
	}

	return userData, nil
}

func FetchName(c *fiber.Ctx, uid interface{}, role interface{}, name string) (*users_get.UserData, error) {
	valid, err := validator.Var(c, name, "user_name")
	if err != nil || !valid {
		return nil, err
	}

	allData := false
	if role != nil {
		allData = utils.In(role.(sqlc.UserRole), []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
	}
	userData, err := GetUserByName(c.Context(), name, uid, allData)
	if err != nil {
		return nil, utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if userData == nil {
		return nil, utils.Error(c, fiber.StatusNotFound, consts.UserNotFound)
	}

	return userData, nil
}

func Route(c *fiber.Ctx) error {
	uid := c.Locals("uid")
	role := c.Locals("role")

	userName := c.Query("name")
	userEmail := c.Query("email")

	if userEmail == "" && userName == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}

	var userData *users_get.UserData
	var err error

	if userEmail != "" {
		userData, err = FetchEmail(c, role, userEmail)
		if err != nil || userData == nil {
			return err
		}
	} else if userName != "" {
		userData, err = FetchName(c, uid, role, userName)
		if err != nil || userData == nil {
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(userData)
}
