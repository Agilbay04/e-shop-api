package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
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
		case "max":
			return fmt.Sprintf("Maximum value is %s", param)
		case "gt":
			return fmt.Sprintf("Value must be greater than %s", param)
		case "gte":
			return fmt.Sprintf("Value must be greater than or equal to %s", param)
		case "oneof":
			return fmt.Sprintf("Value must be one of %s", param)
		case "eqfield":
			return fmt.Sprintf("Value must be equal to %s", param)
	}
	return "Invalid value"
}

// Register json tag name
func RegisterJSONTagName() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
