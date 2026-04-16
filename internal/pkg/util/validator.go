package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Map error validator to map[string]string
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	
	// Check if err is implement validator.ValidationErrors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, f := range validationErrors {
			// extract error message
			errors[f.Field()] = MsgForTag(f.Tag(), f.Param())
		}
	}
	return errors
}

// Message for tag validator
func MsgForTag(tag string, param string) string {
	switch tag {
		case "required":
			return "This field is required"
		case "email":
			return "Invalid email format"
		case "min":
			return fmt.Sprintf("Minimum value is %s", param)
		case "gt":
			return fmt.Sprintf("Value must be greater than %s", param)
		case "gte":
			return fmt.Sprintf("Value must be greater than or equal to %s", param)
		case "oneof":
			return fmt.Sprintf("Value must be one of %s", param)
	}
	return "Invalid value"
}