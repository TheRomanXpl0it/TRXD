package validator

import (
	"encoding/json"
	"strconv"

	"github.com/go-playground/validator/v10"
)

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
