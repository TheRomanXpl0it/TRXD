package attachments_create

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"
	"trxd/validator"

	"github.com/go-playground/form/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func Route(c *fiber.Ctx) error {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidMultipartForm)
	}

	var data struct {
		ChallID *int32 `form:"chall_id" validate:"required,id"`
	}
	decoder := form.NewDecoder()
	if err = decoder.Decode(&data, multipartForm.Value); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFormData)
	}

	if len(multipartForm.File) == 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}

	names := make([]string, 0)
	headers := make([]*multipart.FileHeader, 0)
	for _, files := range multipartForm.File {
		for _, file := range files {
			names = append(names, file.Filename)
			headers = append(headers, file)
		}
	}

	valid, err = validator.Var(c, names, "attachments")
	if err != nil || !valid {
		return err
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	dir := fmt.Sprintf("attachments/%d/", *data.ChallID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingAttachmentsDir, err)
		}
	}

	hashes := make([]string, 0, len(headers))
	for _, file := range headers {
		cleanPath := filepath.Clean(dir + filepath.Base(file.Filename))
		if !strings.HasPrefix(cleanPath, dir) {
			return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFilePath)
		}

		f, err := file.Open()
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorHashingFile, err)
		}
		defer f.Close()

		hash, err := crypto_utils.HashFile(f)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorHashingFile, err)
		}

		hashes = append(hashes, hash)

		hashedPath := dir + hash + "/"
		if _, err := os.Stat(hashedPath); os.IsNotExist(err) {
			err := os.MkdirAll(hashedPath, 0755)
			if err != nil {
				return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingAttachmentsDir, err)
			}
		}

		cleanPath = filepath.Clean(hashedPath + filepath.Base(file.Filename))
		if !strings.HasPrefix(cleanPath, hashedPath) {
			return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFilePath)
		}

		err = c.SaveFile(file, cleanPath)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingFile, err)
		}
	}

	err = CreateAttachments(c.Context(), *data.ChallID, names, hashes)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return utils.Error(c, fiber.StatusConflict, consts.AttachmentAlreadyExists)
			}
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingAttachments, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
