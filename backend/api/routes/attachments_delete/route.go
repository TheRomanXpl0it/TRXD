package attachments_delete

import (
	"fmt"
	"os"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32    `json:"chall_id" validate:"required,id"`
		Names   *[]string `json:"names" validate:"required,attachments"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	challID := *data.ChallID

	hashes := make([]string, 0, len(*data.Names))
	for _, name := range *data.Names {
		hash, err := GetAttachmentHash(c.Context(), challID, name)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingAttachment, err)
		}
		if hash == "" {
			return utils.Error(c, fiber.StatusNotFound, consts.AttachmentNotFound)
		}

		hashes = append(hashes, hash)
	}
	for i, name := range *data.Names {
		err = DeleteAttachment(c.Context(), challID, name)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingAttachment, err)
		}

		err = os.RemoveAll(fmt.Sprintf("attachments/%d/%s", challID, hashes[i]))
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingAttachment, err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
