package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(input interface{}) map[string]string {
	validate := validator.New()

	err := validate.Struct(input)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return map[string]string{"error": "validasi input gagal: " + err.Error()}
		}

		validationErrors := make(map[string]string)
		inputVal := reflect.ValueOf(input)
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := inputVal.Type().FieldByName(err.Field())
			formTag := field.Tag.Get("form")
			jsonTag := field.Tag.Get("json")

			var fieldName string
			if formTag != "" {
				fieldName = formTag
			} else if jsonTag != "" {
				fieldName = jsonTag
			} else {
				fieldName = err.Field()
			}

			validationErrors[fieldName] = getErrorMessage(err)
		}
		return validationErrors
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "oneof":
		fields := strings.Split(err.Param(), " ")
		return fmt.Sprintf("at least one of these fields must be provided: %s", strings.Join(fields, ", "))
	default:
		return fmt.Sprintf("%s is not valid", err.Field())
	}
}
