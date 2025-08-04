package routes

import (
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func CreateChallenge(c *fiber.Ctx) error {
	var data struct {
		Name        string        `json:"name"`
		Category    string        `json:"category"`
		Description string        `json:"description"`
		Type        db.DeployType `json:"type"`
		MaxPoints   int32         `json:"max_points"`
		ScoreType   db.ScoreType  `json:"score_type"`
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
	if !utils.In(data.Type, []db.DeployType{db.DeployTypeNormal, db.DeployTypeContainer, db.DeployTypeCompose}) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallType)
	}
	if data.MaxPoints <= 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallMaxPoints)
	}
	if !utils.In(data.ScoreType, []db.ScoreType{db.ScoreTypeStatic, db.ScoreTypeDynamic}) {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallScoreType)
	}

	challenge, err := db.CreateChallenge(c.Context(), data.Name, data.Category, data.Description, data.Type, data.MaxPoints, data.ScoreType)
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

func CreateFlag(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
		Flag    string `json:"flag"`
		Regex   bool   `json:"regex"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.Flag == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Flag) > consts.MaxFlagLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongFlag)
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	flag, err := db.CreateFlag(c.Context(), *data.ChallID, data.Flag, data.Regex)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorCreatingFlag, err)
	}
	if flag == nil {
		return utils.Error(c, fiber.StatusConflict, consts.FlagAlreadyExists)
	}

	return c.SendStatus(fiber.StatusOK)
}

func GetChallenges(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	role := c.Locals("role").(db.UserRole)

	all := utils.In(role, []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	challenges, err := db.GetChallenges(c.Context(), uid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}

	return c.Status(fiber.StatusOK).JSON(challenges)
}

func GetChallenge(c *fiber.Ctx) error {
	uid := c.Locals("uid").(int32)
	role := c.Locals("role").(db.UserRole)

	challengeID, err := c.ParamsInt("id")
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallengeID)
	}

	all := utils.In(role, []db.UserRole{db.UserRoleAuthor, db.UserRoleAdmin})
	challenge, err := db.GetChallenge(c.Context(), int32(challengeID), uid, all)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenges, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(challenge)
}

func Submit(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
		Flag    string `json:"flag"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.Flag == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Flag) > consts.MaxFlagLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongFlag)
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	uid := c.Locals("uid").(int32)

	status, err := db.SubmitFlag(c.Context(), uid, *data.ChallID, data.Flag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorSubmittingFlag, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": status,
	})
}

func DeleteChallenge(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if *data.ChallID < 0 {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidChallengeID)
	}

	err := db.DeleteChallenge(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingChallenge, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteFlag(c *fiber.Ctx) error {
	var data struct {
		ChallID *int32 `json:"chall_id"`
		Flag    string `json:"flag"`
	}
	if err := c.BodyParser(&data); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, consts.InvalidJSON)
	}

	if data.ChallID == nil || data.Flag == "" {
		return utils.Error(c, fiber.StatusBadRequest, consts.MissingRequiredFields)
	}
	if len(data.Flag) > consts.MaxFlagLength {
		return utils.Error(c, fiber.StatusBadRequest, consts.LongFlag)
	}

	challenge, err := db.GetChallengeByID(c.Context(), *data.ChallID)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorFetchingChallenge, err)
	}
	if challenge == nil {
		return utils.Error(c, fiber.StatusNotFound, consts.ChallengeNotFound)
	}

	err = db.DeleteFlag(c.Context(), *data.ChallID, data.Flag)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, consts.ErrorDeletingFlag, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
