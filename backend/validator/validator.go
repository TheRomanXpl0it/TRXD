package validator

import (
	"fmt"
	"math"
	"strings"
	"trxd/utils"
	"trxd/utils/consts"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate
var uni *ut.UniversalTranslator
var trans ut.Translator

func init() {
	validate = validator.New()

	initTranslation()

	validate.RegisterAlias("id", fmt.Sprintf("min=0,max=%d", math.MaxInt32))
	validate.RegisterAlias("password", fmt.Sprintf("min=%d,max=%d", consts.MinPasswordLen, consts.MaxPasswordLen))
	validate.RegisterValidation("country", validCountry)

	validate.RegisterAlias("category_name", fmt.Sprintf("max=%d", consts.MaxCategoryLen))

	validate.RegisterAlias("challenge_name", fmt.Sprintf("max=%d", consts.MaxChallNameLen))
	validate.RegisterAlias("challenge_description", fmt.Sprintf("max=%d", consts.MaxChallDescLen))
	validate.RegisterAlias("challenge_difficulty", fmt.Sprintf("max=%d", consts.MaxChallDifficultyLen))
	validate.RegisterAlias("challenge_type", "oneof="+strings.Join(consts.DeployTypesStr, " "))
	validate.RegisterAlias("challenge_max_points", fmt.Sprintf("min=0,max=%d", math.MaxInt32))
	validate.RegisterAlias("challenge_score_type", "oneof="+strings.Join(consts.ScoreTypesStr, " "))
	validate.RegisterAlias("challenge_port", fmt.Sprintf("min=%d,max=%d", consts.MinPort, consts.MaxPort))
	validate.RegisterAlias("challenge_lifetime", fmt.Sprintf("min=0,max=%d", math.MaxInt32))
	validate.RegisterValidation("challenge_envs", validJson)
	validate.RegisterAlias("challenge_max_memory", fmt.Sprintf("min=0,max=%d", math.MaxInt32))
	validate.RegisterValidation("challenge_max_cpu", validFloat)

	validate.RegisterAlias("flag", fmt.Sprintf("max=%d", consts.MaxFlagLen))

	validate.RegisterAlias("tag_name", fmt.Sprintf("max=%d", consts.MaxTagNameLen))

	validate.RegisterAlias("team_name", fmt.Sprintf("max=%d", consts.MaxTeamNameLen))

	validate.RegisterAlias("user_name", fmt.Sprintf("max=%d", consts.MaxUserNameLen))
	validate.RegisterAlias("user_email", fmt.Sprintf("max=%d,email", consts.MaxEmailLen))
}

func errHandle(c *fiber.Ctx, err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InternalServerError, err)
	}

	errs := err.(validator.ValidationErrors)
	if len(errs) == 0 {
		return utils.Error(c, fiber.StatusInternalServerError, consts.InternalServerError, err)
	}

	return utils.Error(c, fiber.StatusBadRequest, errs[0].Translate(trans))
}

func Struct(c *fiber.Ctx, s interface{}) (bool, error) {
	err := validate.Struct(s)
	if err != nil {
		return false, errHandle(c, err)
	}

	return true, nil
}

func Var(c *fiber.Ctx, v interface{}, tag string) (bool, error) {
	err := validate.Var(v, tag)
	if err != nil {
		return false, errHandle(c, err)
	}

	return true, nil
}
