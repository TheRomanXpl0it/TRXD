package users_info

import (
	"fmt"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	uidLocal := c.Locals("uid")
	if uidLocal == nil {
		return c.SendStatus(fiber.StatusOK)
	}

	uid := uidLocal.(int32)

	user, err := db.GetUserByID(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, fmt.Errorf("user not found"))
	}

	var teamID *int32
	if user.TeamID.Valid {
		teamID = &user.TeamID.Int32
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      user.ID,
		"name":    user.Name,
		"role":    user.Role,
		"team_id": teamID,
	})
}
