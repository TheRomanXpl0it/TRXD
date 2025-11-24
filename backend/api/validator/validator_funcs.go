package validator

import (
	"math"
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
	// TODO: make this a map[string]string
	return err == nil
}

func validFloat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	res, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}

	return 0.0 < res && res <= math.MaxInt32
}
