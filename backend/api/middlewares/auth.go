package middlewares

import (
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func withUser(c *fiber.Ctx, requireAuth bool, allowedRoles []sqlc.UserRole) error {
	sess, err := db.Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		if !requireAuth {
			return c.Next()
		}
		return utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	user, err := db.GetUserByID(c.Context(), uid.(int32))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil || !utils.In(user.Role, allowedRoles) {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}

	tid := int32(-1)
	if user.TeamID.Valid {
		tid = user.TeamID.Int32
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)
	c.Locals("tid", tid)

	return c.Next()
}

func NoAuth(c *fiber.Ctx) error {
	return withUser(c, false, []sqlc.UserRole{sqlc.UserRoleSpectator, sqlc.UserRolePlayer, sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
}

func Spectator(c *fiber.Ctx) error {
	return withUser(c, true, []sqlc.UserRole{sqlc.UserRoleSpectator, sqlc.UserRolePlayer, sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
}

func Player(c *fiber.Ctx) error {
	return withUser(c, true, []sqlc.UserRole{sqlc.UserRolePlayer, sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
}

func Author(c *fiber.Ctx) error {
	return withUser(c, true, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})
}

func Admin(c *fiber.Ctx) error {
	return withUser(c, true, []sqlc.UserRole{sqlc.UserRoleAdmin})
}

func Team(c *fiber.Ctx) error {
	tid := c.Locals("tid")
	role := c.Locals("role")

	if role == nil ||
		(role.(sqlc.UserRole) == sqlc.UserRolePlayer &&
			(tid == nil || tid.(int32) == -1)) {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}

	return c.Next()
}
