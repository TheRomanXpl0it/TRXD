package validator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var validate *validator.Validate
var uni *ut.UniversalTranslator
var trans ut.Translator

func init() {
	validate = validator.New()

	initTranslation()

	deployTypes := []string{string(sqlc.DeployTypeNormal), string(sqlc.DeployTypeContainer), string(sqlc.DeployTypeCompose)}
	scoreTypes := []string{string(sqlc.ScoreTypeStatic), string(sqlc.ScoreTypeDynamic)}

	// TODO: put a max to integers (ex: "challenge_max_points")

	validate.RegisterAlias("category_name", fmt.Sprintf("max=%d", consts.MaxCategoryLength))
	validate.RegisterAlias("category_icon", fmt.Sprintf("max=%d", consts.MaxIconLength))

	validate.RegisterAlias("challenge_id", "min=0")
	validate.RegisterAlias("challenge_name", fmt.Sprintf("max=%d", consts.MaxChallNameLength))
	validate.RegisterAlias("challenge_description", fmt.Sprintf("max=%d", consts.MaxChallDescLength))
	validate.RegisterAlias("challenge_difficulty", fmt.Sprintf("max=%d", consts.MaxChallDifficultyLength))
	validate.RegisterAlias("challenge_type", "oneof="+strings.Join(deployTypes, " "))
	validate.RegisterAlias("challenge_max_points", "min=0")
	validate.RegisterAlias("challenge_score_type", "oneof="+strings.Join(scoreTypes, " "))
	validate.RegisterAlias("challenge_port", fmt.Sprintf("min=%d,max=%d", consts.MinPort, consts.MaxPort))
	validate.RegisterAlias("challenge_lifetime", "min=0")
	validate.RegisterValidation("challenge_envs", validJson)
	validate.RegisterAlias("challenge_max_memory", "min=0")
	validate.RegisterValidation("challenge_max_cpu", validFloat)

	validate.RegisterAlias("flag_flag", fmt.Sprintf("max=%d", consts.MaxFlagLength))

	validate.RegisterAlias("tag_name", fmt.Sprintf("max=%d", consts.MaxTagNameLength))

	validate.RegisterAlias("team_id", "min=0")
	validate.RegisterAlias("team_name", fmt.Sprintf("max=%d", consts.MaxNameLength))
	validate.RegisterAlias("team_password", fmt.Sprintf("min=%d,max=%d", consts.MinPasswordLength, consts.MaxPasswordLength))
	validate.RegisterAlias("team_country", fmt.Sprintf("max=%d", consts.MaxCountryLength))
	validate.RegisterAlias("team_image", fmt.Sprintf("max=%d", consts.MaxImageLength))
	validate.RegisterAlias("team_bio", fmt.Sprintf("max=%d", consts.MaxBioLength))

	validate.RegisterAlias("user_id", "min=0")
	validate.RegisterAlias("user_name", fmt.Sprintf("max=%d", consts.MaxNameLength))
	validate.RegisterAlias("user_email", fmt.Sprintf("max=%d", consts.MaxEmailLength))
	validate.RegisterAlias("user_password", fmt.Sprintf("min=%d,max=%d", consts.MinPasswordLength, consts.MaxPasswordLength))
	validate.RegisterAlias("user_country", fmt.Sprintf("max=%d", consts.MaxCountryLength))
	validate.RegisterAlias("user_image", fmt.Sprintf("max=%d", consts.MaxImageLength))
}

func validJson(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	var tmp map[string]string
	err := json.Unmarshal([]byte(value), &tmp)
	return err == nil
}

func validFloat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	_, err := strconv.ParseFloat(value, 32)
	return err == nil
}

func initTranslation() {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	registerTranslation("required", consts.MissingRequiredFields)
	registerTranslation("min", consts.MinError)
	registerTranslation("max", consts.MaxError)
	registerTranslation("oneof", consts.OneofError)
	registerTranslation("challenge_envs", consts.InvalidEnvs)
	registerTranslation("challenge_max_cpu", consts.InvalidMaxCpu)
}

func registerTranslation(tag string, format string) {
	err := validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, format, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		var t string
		field := fe.Field()
		if field != "" {
			t, _ = ut.T(tag, field, fe.Param())
		} else {
			t, _ = ut.T(tag, fe.Tag(), fe.Param())
		}
		return t
	})
	if err != nil {
		log.Error("Error Registering Translation", "tag", tag, "err", err)
	}
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
