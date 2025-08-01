package auth

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// TODO: set store as redis + set configs (expire time, etc.)
var Store = session.New(session.Config{})

// TODO: make auth_test.go

func AuthRequired(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingSession, err)
	}

	uid := sess.Get("uid")
	if uid == nil {
		return utils.Error(c, fiber.StatusUnauthorized, consts.Unauthorized)
	}

	c.Locals("uid", uid)
	return c.Next()
}

func PlayerRequired(c *fiber.Ctx) error {
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
	if user == nil || (user.Role != db.UserRolePlayer &&
		user.Role != db.UserRoleAuthor && user.Role != db.UserRoleAdmin) {
		return utils.Error(c, fiber.StatusForbidden, consts.Unauthorized)
	}

	if user.Role == db.UserRolePlayer {
		if user.TeamID.Valid {
			c.Locals("tid", user.TeamID.Int32)
		} else {
			c.Locals("tid", int32(-1))
		}
	}
	c.Locals("uid", uid)
	return c.Next()
}

func AuthorRequired(c *fiber.Ctx) error {
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
	if user == nil || (user.Role != db.UserRoleAuthor && user.Role != db.UserRoleAdmin) {
		return utils.Error(c, fiber.StatusForbidden, consts.Unauthorized)
	}

	c.Locals("uid", uid)
	return c.Next()
}

func AdminRequired(c *fiber.Ctx) error {
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
	if user == nil || user.Role != db.UserRoleAdmin {
		return utils.Error(c, fiber.StatusForbidden, consts.Unauthorized)
	}

	c.Locals("uid", uid)
	return c.Next()
}
