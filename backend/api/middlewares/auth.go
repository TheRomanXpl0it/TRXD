package middlewares

import (
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
)

func NoAuth(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return c.Next()
	}

	user, err := db.GetUserByID(c.Context(), uid.(int32))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)

	return c.Next()
}

func Spectator(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	user, err := db.GetUserByID(c.Context(), uid.(int32))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)

	return c.Next()
}

func Team(c *fiber.Ctx) error {
	tid := c.Locals("tid").(int32)
	role := c.Locals("role").(sqlc.UserRole)

	if role == sqlc.UserRolePlayer && tid == -1 {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}

	return c.Next()
}

func Player(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	user, err := db.GetUserByID(c.Context(), uid.(int32))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil || (user.Role != sqlc.UserRolePlayer &&
		user.Role != sqlc.UserRoleAuthor && user.Role != sqlc.UserRoleAdmin) {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)

	return c.Next()
}

func Author(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	user, err := db.GetUserByID(c.Context(), uid.(int32))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil || (user.Role != sqlc.UserRoleAuthor && user.Role != sqlc.UserRoleAdmin) {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)
	return c.Next()
}

func Admin(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	user, err := db.GetUserByID(c.Context(), uid.(int32))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingUser, err)
	}
	if user == nil || user.Role != sqlc.UserRoleAdmin {
		return utils.Error(c, fiber.StatusForbidden, consts.Forbidden)
	}

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)
	return c.Next()
}
