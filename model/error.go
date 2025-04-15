package model

import (
	"fmt"
	"net/http"
)

// ErrorResponse represents a standardized error response format
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response with the given status code and message
func NewErrorResponse(statusCode int, number int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    formatErrorCode(statusCode, number),
		Message: message,
	}
}

// Helper function to format error code in "[status code]-[number]" format
func formatErrorCode(statusCode int, number int) string {
	return fmt.Sprintf("%d-%d", statusCode, number)
}

// Auth error constants
var (
	// 400 Bad Request errors
	InvalidRequestBodyResponse = NewErrorResponse(http.StatusBadRequest, 1, "Invalid request body")
	ValidationFailedResponse   = NewErrorResponse(http.StatusBadRequest, 2, "Validation failed")

	// 401 Unauthorized errors
	InvalidCredentialsResponse      = NewErrorResponse(http.StatusUnauthorized, 1, "Invalid email or password")
	MissingAuthHeaderResponse       = NewErrorResponse(http.StatusUnauthorized, 2, "Missing authorization header")
	InvalidAuthHeaderFormatResponse = NewErrorResponse(http.StatusUnauthorized, 3, "Invalid authorization header format")
	TokenExpiredResponse            = NewErrorResponse(http.StatusUnauthorized, 4, "Token expired")
	InvalidTokenResponse            = NewErrorResponse(http.StatusUnauthorized, 5, "Invalid token")
	InvalidTokenTypeResponse        = NewErrorResponse(http.StatusUnauthorized, 6, "Invalid token type")

	// 409 Conflict errors
	EmailAlreadyExistsResponse = NewErrorResponse(http.StatusConflict, 1, "Email already exists")

	// 500 Internal Server Error errors
	FailedToCreateUserResponse    = NewErrorResponse(http.StatusInternalServerError, 1, "Failed to create user")
	AuthenticationFailedResponse  = NewErrorResponse(http.StatusInternalServerError, 2, "Authentication failed")
	FailedToGetUserClaimsResponse = NewErrorResponse(http.StatusInternalServerError, 3, "Failed to get user claims")
)
