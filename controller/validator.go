package controller

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator is the validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the request against the struct validation tags
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewValidator creates a new validator for Echo
func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}
