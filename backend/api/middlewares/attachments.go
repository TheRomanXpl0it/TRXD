package middlewares

import (
	"strconv"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Attachments(c *fiber.Ctx) error {
	path := strings.Split(c.Path(), "/")
	if len(path) != 4 || len(path[2]) == 0 || len(path[3]) == 0 {
		return utils.Error(c, fiber.StatusNotFound, consts.NotFound)
	}

	challID, err := strconv.Atoi(path[2])
	if err != nil {
		return utils.Error(c, fiber.StatusNotFound, consts.NotFound)
	}

	valid, err := validator.Var(c, challID, "id")
	if err != nil || !valid {
		return utils.Error(c, fiber.StatusNotFound, consts.NotFound)
	}

	role := c.Locals("role").(sqlc.UserRole)

	res, err := db.GetHiddenAndAttachments(c.Context(), int32(challID))
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InternalServerError, err)
	}
	if res == nil ||
		(res.Hidden && !utils.In(role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})) ||
		!utils.In(c.Path()[1:], strings.Split(res.Attachments, consts.Separator)) {
		return utils.Error(c, fiber.StatusNotFound, consts.NotFound)
	}

	return c.Next()
}
