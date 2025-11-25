package challenges_create

import (
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func Route(c *fiber.Ctx) error {
	var data struct {
		Name        string          `json:"name" validate:"required,challenge_name"`
		Category    string          `json:"category" validate:"required,category_name"`
		Description string          `json:"description" validate:"required,challenge_description"`
		Type        sqlc.DeployType `json:"type" validate:"required,challenge_type"`
		MaxPoints   int32           `json:"max_points" validate:"required,challenge_max_points"`
		ScoreType   sqlc.ScoreType  `json:"score_type" validate:"required,challenge_score_type"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	valid, err := validator.Struct(c, data)
	if err != nil || !valid {
		return err
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

	return c.SendStatus(fiber.StatusOK)
}
