package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ValidatePayload(payload interface{}) (map[string]string, error) {
	err := Validate.Struct(payload)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			invalidFields := make(map[string]string)
			for _, e := range validationErrors {
				invalidFields[e.Field()] = fmt.Sprintf("Validation failed on the '%s' tag", e.Tag())
			}
			return invalidFields, fmt.Errorf("invalid payload")
		}
	}
	return nil, nil
}
