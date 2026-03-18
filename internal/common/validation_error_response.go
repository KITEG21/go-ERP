package common

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func (v ValidationErrorResponse) CustomErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("It has to have min %s characters", err.Param())
	case "max":
		return fmt.Sprintf("It can not exceed %s characters", err.Param())
	case "email":
		return "It has to be a valid email address"
	case "eqfield":
		return "Passwords do not match"
	case "gt":
		return fmt.Sprintf("It has to be greater than %s", err.Param())
	case "gte":
		if err.Param() == "0" {
			return "It has to be a positive number"
		}
		return fmt.Sprintf("It has to be greater than or equal to %s", err.Param())
	default:
		return fmt.Sprintf("Validation '%s' failed", err.Tag())
	}
}
