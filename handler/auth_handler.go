// Package handler contains handler functions
package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
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
