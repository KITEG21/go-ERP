package common

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("notblank", validators.NotBlank)
	return v
}

func ValidateStruct(v *validator.Validate, s interface{}) ([]ValidationErrorResponse, error) {
	err := v.Struct(s)
	if err == nil {
		return nil, nil
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		resp := make([]ValidationErrorResponse, 0, len(ve))
		for _, fe := range ve {
			resp = append(resp, ValidationErrorResponse{
				Field:   fe.Field(),
				Tag:     fe.Tag(),
				Value:   fe.Param(),
				Message: ValidationErrorResponse{}.CustomErrorMessage(fe),
			})
		}
		return resp, err
	}
	return nil, err
}
