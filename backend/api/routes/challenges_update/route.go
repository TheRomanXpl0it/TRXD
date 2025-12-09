package challenges_update

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/go-playground/form/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type UpdateChallParams struct {
	ChallID     *int32           `form:"chall_id" validate:"required,id"`
	Name        string           `form:"name" validate:"challenge_name"`
	Category    string           `form:"category" validate:"category_name"`
	Description *string          `form:"description" validate:"omitempty,challenge_description"`
	Difficulty  *string          `form:"difficulty" validate:"omitempty,challenge_difficulty"`
	Authors     *[]string        `form:"authors"` // TODO: add a valdator for authors
	Type        *sqlc.DeployType `form:"type" validate:"omitempty,challenge_type"`
	Hidden      *bool            `form:"hidden"`
	MaxPoints   *int32           `form:"max_points" validate:"omitempty,challenge_max_points"`
	ScoreType   *sqlc.ScoreType  `form:"score_type" validate:"omitempty,challenge_score_type"`
	Host        *string          `form:"host"`
	Port        *int32           `form:"port" validate:"omitempty,challenge_port"`
	Attachments *[]string        `form:"attachments"`

	Image      *string `form:"image"`
	Compose    *string `form:"compose"`
	HashDomain *bool   `form:"hash_domain"`
	Lifetime   *int32  `form:"lifetime" validate:"omitempty,challenge_lifetime"`
	Envs       *string `form:"envs" validate:"omitempty,challenge_envs"`
	MaxMemory  *int32  `form:"max_memory" validate:"omitempty,challenge_max_memory"`
	MaxCpu     *string `form:"max_cpu" validate:"omitempty,challenge_max_cpu"`
}

func Route(c *fiber.Ctx) error {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidMultipartForm)
	}

	data := &UpdateChallParams{}
	decoder := form.NewDecoder()
	err = decoder.Decode(&data, multipartForm.Value)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFormData)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
	}
	if IsChallEmpty(data) && IsDockerConfigsEmpty(data) && len(multipartForm.File) == 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.NoDataToUpdate)
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	attachments := make(map[string]*multipart.FileHeader)
	dir := fmt.Sprintf("attachments/%d/", *data.ChallID)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingAttachmentsDir, err)
		}
	}

	attachmentsList := make([]string, 0)
	for _, files := range multipartForm.File {
		for _, file := range files {
			cleanPath := filepath.Clean(dir + filepath.Base(file.Filename))
			if !strings.HasPrefix(cleanPath, dir) {
				return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFilePath)
			}
			attachments[cleanPath] = file
			attachmentsList = append(attachmentsList, cleanPath)
		}
	}

	if len(attachmentsList) != 0 {
		data.Attachments = &attachmentsList
	}

	err = UpdateChallenge(c.Context(), data)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return utils.Error(c, fiber.StatusConflict, consts.ChallengeNameAlreadyExists)
			}
			if pqErr.Code == "23503" { // Foreign key violation error code
				return utils.Error(c, fiber.StatusNotFound, consts.CategoryNotFound)
			}
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingChallenge, err)
	}

	for cleanPath, file := range attachments {
		err := c.SaveFile(file, cleanPath)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSavingFile, err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
