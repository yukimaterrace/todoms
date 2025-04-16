package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yukimaterrace/todoms/handler"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
)

// AuthController handles authentication related HTTP requests
type AuthController struct {
	authService service.AuthenticationService
	userService service.UserService
	authHandler *handler.AuthHandler
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService service.AuthenticationService, userService service.UserService, authHandler *handler.AuthHandler) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
		authHandler: authHandler,
	}
}

// RegisterRoutes registers the auth routes to the given Echo instance
func (c *AuthController) RegisterRoutes(e *echo.Echo) {
	auth := e.Group("/api/auth")
	auth.POST("/signup", c.SignUp)
	auth.POST("/login", c.Login)
	auth.POST("/refresh", c.Refresh)
	auth.GET("/me", c.Me, c.authHandler.RequireAuth)
}

// SignUp handles user registration
func (c *AuthController) SignUp(ctx echo.Context) error {
	req := new(model.SignUpRequest)
	if err := ValidateRequest(ctx, req); err != nil {
		return err // Error response already sent by ValidateRequest
	}

	user, err := c.userService.CreateUser(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			return ctx.JSON(http.StatusConflict, model.EmailAlreadyExistsResponse)
		}
		return ctx.JSON(http.StatusInternalServerError, model.FailedToCreateUserResponse)
	}

	return ctx.JSON(http.StatusCreated, model.UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	})
}

// Login handles user authentication and returns JWT tokens
func (c *AuthController) Login(ctx echo.Context) error {
	req := new(model.LoginRequest)
	if err := ValidateRequest(ctx, req); err != nil {
		return err // Error response already sent by ValidateRequest
	}

	tokenPair, err := c.authService.Authenticate(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrInvalidCredentials:
			return ctx.JSON(http.StatusUnauthorized, model.InvalidCredentialsResponse)
		default:
			return ctx.JSON(http.StatusInternalServerError, model.AuthenticationFailedResponse)
		}
	}

	return ctx.JSON(http.StatusOK, tokenPair)
}

// Refresh handles token refresh and returns a new token pair
func (c *AuthController) Refresh(ctx echo.Context) error {
	req := new(model.RefreshTokenRequest)
	if err := ValidateRequest(ctx, req); err != nil {
		return err // Error response already sent by ValidateRequest
	}

	tokenPair, err := c.authService.RefreshToken(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidToken, service.ErrExpiredToken:
			return ctx.JSON(http.StatusUnauthorized, model.InvalidTokenResponse)
		case service.ErrInvalidTokenType:
			return ctx.JSON(http.StatusUnauthorized, model.InvalidTokenTypeResponse)
		case service.ErrUserNotFound:
			return ctx.JSON(http.StatusUnauthorized, model.InvalidCredentialsResponse)
		default:
			return ctx.JSON(http.StatusInternalServerError, model.AuthenticationFailedResponse)
		}
	}

	return ctx.JSON(http.StatusOK, tokenPair)
}

// Me returns information about the authenticated user
func (c *AuthController) Me(ctx echo.Context) error {
	claims, err := c.authHandler.GetUserClaims(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.FailedToGetUserClaimsResponse)
	}

	return ctx.JSON(http.StatusOK, model.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
	})
}
