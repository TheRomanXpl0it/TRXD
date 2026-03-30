package users_info

import (
	"fmt"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	emailVerification, err := db.GetConfig(c.Context(), "email-verification")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	startTime, err := db.GetConfig(c.Context(), "start-time")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	endTime, err := db.GetConfig(c.Context(), "end-time")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}

	info := fiber.Map{
		"email_verification": emailVerification == "true",
	}
	if startTime != "" {
		info["start_time"] = startTime
	}
	if endTime != "" {
		info["end_time"] = endTime
	}

	uidLocal := c.Locals("uid")
	if uidLocal == nil {
		return c.Status(fiber.StatusOK).JSON(info)
	}

	uid := uidLocal.(int32)

	user, err := db.GetUserByID(c.Context(), uid)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, fmt.Errorf("user not found"))
	}
	info["id"] = user.ID
	info["name"] = user.Name
	info["role"] = user.Role

	userMode, err := db.GetConfig(c.Context(), "user-mode")
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingConfig, err)
	}
	info["user_mode"] = userMode == "true"

	var teamID *int32
	if user.TeamID.Valid {
		teamID = &user.TeamID.Int32
	}
	info["team_id"] = teamID

	return c.Status(fiber.StatusOK).JSON(info)
}
