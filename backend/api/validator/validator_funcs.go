package validator

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

func validCountry(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	err := validate.Var(value, "iso3166_1_alpha3")
	return err == nil
}

func validHttpUrl(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	err := validate.Var(value, "http_url")
	return err == nil
}

func validJson(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	err := validate.Var(value, "json")
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
