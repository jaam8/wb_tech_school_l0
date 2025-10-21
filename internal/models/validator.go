package models

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("alphaunicode_with_space", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		if val == "" {
			return false
		}
		for _, r := range val {
			if !(unicode.IsLetter(r) || unicode.IsSpace(r) || r == '.' || r == '-') {
				return false
			}
		}

		return true
	})
}
