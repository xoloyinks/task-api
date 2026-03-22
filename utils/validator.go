package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = formatMessage(e)
	}

	return errors
}

func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s cannot exceed %s characters", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
