package validator

import (
	"trxd/utils/consts"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/tde-nico/log"
)

func initTranslation() {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	registerTranslation("required", consts.MissingRequiredFields)
	registerTranslation("min", consts.MinError)
	registerTranslation("max", consts.MaxError)
	registerTranslation("oneof", consts.OneofError)
	registerTranslation("email", consts.InvalidEmail)

	registerTranslation("image", consts.InvalidHttpUrl)
	registerTranslation("country", consts.InvalidCountry)
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
