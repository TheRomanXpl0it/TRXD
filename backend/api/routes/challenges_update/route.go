package challenges_update

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/go-playground/form/v4"
	"github.com/gofiber/fiber/v2"
)

type UpdateChallParams struct {
	ChallID     *int32           `form:"chall_id"`
	Name        string           `form:"name"`
	Category    string           `form:"category"`
	Description string           `form:"description"`
	Difficulty  string           `form:"difficulty"`
	Authors     []string         `form:"authors"`
	Type        *sqlc.DeployType `form:"type"`
	Hidden      *bool            `form:"hidden"`
	MaxPoints   *int             `form:"max_points"`
	ScoreType   *sqlc.ScoreType  `form:"score_type"`
	Host        string           `form:"host"`
	Port        *int             `form:"port"`
	Attachments []string

	Image      string `form:"image"`
	Compose    string `form:"compose"`
	HashDomain *bool  `form:"hash_domain"`
	Lifetime   *int   `form:"lifetime"`
	Envs       string `form:"envs"`
	MaxMemory  *int   `form:"max_memory"`
	MaxCpu     string `form:"max_cpu"`
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

	if data.ChallID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.ChallIDRequired)
	}
	if IsChallEmpty(data) && IsDockerConfigsEmpty(data) && len(multipartForm.File) == 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.NoDataToUpdate)
	}

	if len(data.Name) > consts.MaxChallNameLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongChallName)
	}
	if len(data.Category) > consts.MaxCategoryLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongCategory)
	}
	if len(data.Description) > consts.MaxChallDescLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongChallDesc)
	}
	if len(data.Difficulty) > consts.MaxChallDifficultyLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongChallDifficulty)
	}
	if data.MaxPoints != nil && *data.MaxPoints < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallMaxPoints)
	}
	if data.Port != nil && (*data.Port < consts.MinPort || *data.Port > consts.MaxPort) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidPort)
	}

	if data.Lifetime != nil && *data.Lifetime <= 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidLifetime)
	}
	if data.Envs != "" {
		var tmp map[string]string
		err = json.Unmarshal([]byte(data.Envs), &tmp)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, consts.InvalidEnvs)
		}
	}
	if data.MaxMemory != nil && *data.MaxMemory <= 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidMaxMemory)
	}
	if data.MaxCpu != "" {
		_, err = strconv.ParseFloat(data.MaxCpu, 32)
		if err != nil {
			return utils.Error(c, fiber.StatusBadRequest, consts.InvalidMaxCpu)
		}
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

	for _, files := range multipartForm.File {
		for _, file := range files {
			cleanPath := filepath.Clean(dir + filepath.Base(file.Filename))
			if !strings.HasPrefix(cleanPath, dir) {
				return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFilePath)
			}
			attachments[cleanPath] = file
			data.Attachments = append(data.Attachments, cleanPath)
		}
	}

	err = UpdateChallenge(c.Context(), data)
	if err != nil {
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
