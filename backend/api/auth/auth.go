package auth

import (
	"time"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// TODO: set store as redis
var Store = session.New(session.Config{
	Expiration:     30 * 24 * time.Hour,
	CookiePath:     "/",
	CookieSameSite: fiber.CookieSameSiteLaxMode,
})

func AuthRequired(c *fiber.Ctx) error {
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

// TODO: tests
func TeamRequired(c *fiber.Ctx) error {
	tid := c.Locals("tid").(int32)
	role := c.Locals("role").(db.UserRole)

	if role == db.UserRolePlayer && tid == -1 {
		return utils.Error(c, fiber.StatusForbidden, consts.Unauthorized)
	}

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

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)

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

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)
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

	if user.TeamID.Valid {
		c.Locals("tid", user.TeamID.Int32)
	} else {
		c.Locals("tid", int32(-1))
	}
	c.Locals("uid", uid)
	c.Locals("role", user.Role)
	return c.Next()
}
