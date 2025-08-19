package challenges_create

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name        string          `json:"name"`
		Category    string          `json:"category"`
		Description string          `json:"description"`
		Type        sqlc.DeployType `json:"type"`
		MaxPoints   int32           `json:"max_points"`
		ScoreType   sqlc.ScoreType  `json:"score_type"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.Name == "" || data.Category == "" || data.Description == "" || data.Type == "" || data.ScoreType == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
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
	if !utils.In(data.Type, []sqlc.DeployType{sqlc.DeployTypeNormal, sqlc.DeployTypeContainer, sqlc.DeployTypeCompose}) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallType)
	}
	if data.MaxPoints <= 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallMaxPoints)
	}
	if !utils.In(data.ScoreType, []sqlc.ScoreType{sqlc.ScoreTypeStatic, sqlc.ScoreTypeDynamic}) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallScoreType)
	}

	challenge, err := CreateChallenge(c.Context(), data.Name, data.Category, data.Description, data.Type, data.MaxPoints, data.ScoreType)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" { // Foreign key violation error code
				return utils.Error(c, fiber.StatusNotFound, consts.CategoryNotFound)
			}
		}
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusConflict, consts.ChallengeAlreadyExists)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"name":     challenge.Name,
		"category": challenge.Category,
	})
}
