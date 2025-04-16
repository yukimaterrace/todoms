// Package handler contains handler functions
package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
)

// Error constants for handler
var (
	ErrUserClaimsNotFound  = errors.New("user claims not found in context")
	ErrInvalidUserIDFormat = errors.New("invalid user ID format")
)

// AuthHandler contains authentication handler functions
type AuthHandler struct {
	authService service.AuthenticationService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService service.AuthenticationService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// GetUserClaims retrieves the user claims from the context
func (h *AuthHandler) GetUserClaims(ctx echo.Context) (*service.Claims, error) {
	claims, ok := ctx.Get("user").(*service.Claims)
	if !ok {
		return nil, ErrUserClaimsNotFound
	}
	return claims, nil
}

// GetUserIDFromContext retrieves the user claims and parses the user ID as UUID
// If there's an error, it will return the appropriate HTTP response
func (h *AuthHandler) GetUserIDFromContext(ctx echo.Context) (uuid.UUID, error) {
	claims, err := h.GetUserClaims(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	// Parse user ID from claims
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, ErrInvalidUserIDFormat
	}

	return userID, nil
}

// GetUserIDFromContextWithResponse retrieves user ID and handles error responses
// Returns the user ID and true if successful, or false if an error occurred (in which case the response has already been sent)
func (h *AuthHandler) GetUserIDFromContextWithResponse(ctx echo.Context) (uuid.UUID, bool) {
	userID, err := h.GetUserIDFromContext(ctx)
	if err != nil {
		switch err {
		case ErrUserClaimsNotFound:
			ctx.JSON(http.StatusInternalServerError, model.FailedToGetUserClaimsResponse)
		case ErrInvalidUserIDFormat:
			ctx.JSON(http.StatusInternalServerError, model.InvalidUserIDFormatResponse)
		default:
			ctx.JSON(http.StatusInternalServerError, model.FailedToGetUserClaimsResponse)
		}
		return uuid.Nil, false
	}
	return userID, true
}

// RequireAuth is a middleware to ensure the request is authenticated
func (h *AuthHandler) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return ctx.JSON(http.StatusUnauthorized, model.MissingAuthHeaderResponse)
		}

		// Extract the token from the Authorization header
		// Format: "Bearer {token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			return ctx.JSON(http.StatusUnauthorized, model.InvalidAuthHeaderFormatResponse)
		}

		token := parts[1]
		claims, err := h.authService.ValidateToken(token)
		if err != nil {
			switch err {
			case service.ErrExpiredToken:
				return ctx.JSON(http.StatusUnauthorized, model.TokenExpiredResponse)
			default:
				return ctx.JSON(http.StatusUnauthorized, model.InvalidTokenResponse)
			}
		}

		// Check if it's an access token
		if claims.Type != string(service.AccessToken) {
			return ctx.JSON(http.StatusUnauthorized, model.InvalidTokenTypeResponse)
		}

		// Set the user claims in the context for later use
		ctx.Set("user", claims)

		return next(ctx)
	}
}
