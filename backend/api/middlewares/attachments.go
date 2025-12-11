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
	if len(path) != 5 || // "" + "attachments" + "{chall_id}" + "{file_hash}" + "{file_name}"
		len(path[2]) == 0 || // chall_id
		len(path[3]) == 0 || // file_hash
		len(path[4]) == 0 { // file_name
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
	if res == nil || // challenge not found
		(res.Hidden && !utils.In(role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin})) || // hidden challenge and not author/admin
		!utils.In(path[3]+"/"+path[4], res.Attachments) { // attachment not found
		return utils.Error(c, fiber.StatusNotFound, consts.NotFound)
	}

	return c.Next()
}
