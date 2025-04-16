package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yukimaterrace/todoms/model"
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

// ValidateRequest binds and validates a request, returning appropriate error responses if needed
func ValidateRequest(ctx echo.Context, req interface{}) error {
	// Bind request
	if err := ctx.Bind(req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.InvalidRequestBodyResponse)
		return err
	}

	// Validate request
	if err := ctx.Validate(req); err != nil {
		validationErr := model.ValidationFailedResponse
		validationErr.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, validationErr)
		return err
	}

	return nil
}
