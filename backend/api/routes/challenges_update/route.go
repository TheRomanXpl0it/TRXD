package challenges_update

import (
	"fmt"
	"path/filepath"
	"strings"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/go-playground/form"
	"github.com/gofiber/fiber/v2"
	"github.com/tde-nico/log"
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

	Image      string `form:"image"`
	Compose    string `form:"compose"`
	HashDomain *bool  `form:"hash_domain"`
	Lifetime   *int   `form:"lifetime"`
	Envs       string `form:"envs"`
	MaxMemory  *int   `form:"max_memory"`
	MaxCpu     string `form:"max_cpu"`
}

func Route(c *fiber.Ctx) error {
	multipart, err := c.MultipartForm()
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidMultipartForm, err)
	}

	var data UpdateChallParams
	decoder := form.NewDecoder()
	err = decoder.Decode(&data, multipart.Value)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidFormData, err)
	}

	log.Infof("Decoded form data: %+v", data)

	// TODO: validate all fields

	// TODO: use a transaction

	err = UpdateChallenge(c.Context(), data)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingChallenge, err)
	}

	err = UpdateDockerConfigs(c.Context(), data)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorUpdatingDockerConfigs, err)
	}

	dir := fmt.Sprintf("attachments/%d", *data.ChallID)
	for _, files := range multipart.File {
		for _, file := range files {
			filePath := fmt.Sprintf("%s/%s", dir, file.Filename)
			cleanPath := filepath.Clean(filePath)
			if !strings.HasPrefix(cleanPath, dir) {
				return utils.Error(c, fiber.StatusForbidden, consts.InvalidFilePath)
			}
			c.SaveFile(file, cleanPath)
		}
	}

	// TODO: tests
	log.Critical("DO NOT USE THIS ENDPOINT, IT'S NOT FINISHED NOR TESTED")

	return c.SendStatus(fiber.StatusOK)
}
