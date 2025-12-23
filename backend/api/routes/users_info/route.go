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

	userMode, err := db.GetConfig(c.Context(), "user-mode")
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

	var teamID *int32
	if user.TeamID.Valid {
		teamID = &user.TeamID.Int32
	}
	info := fiber.Map{
		"id":        user.ID,
		"name":      user.Name,
		"role":      user.Role,
		"team_id":   teamID,
		"user_mode": userMode == "true",
	}

	if startTime != "" {
		info["start_time"] = startTime
	}
	if endTime != "" {
		info["end_time"] = endTime
	}

	return c.Status(fiber.StatusOK).JSON(info)
}
