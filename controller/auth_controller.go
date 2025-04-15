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
	auth.GET("/me", c.Me, c.authHandler.RequireAuth)
}

// SignUp handles user registration
func (c *AuthController) SignUp(ctx echo.Context) error {
	req := new(model.SignUpRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.InvalidRequestBodyResponse)
	}

	if err := ctx.Validate(req); err != nil {
		validationErr := model.ValidationFailedResponse
		validationErr.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, validationErr)
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
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, model.InvalidRequestBodyResponse)
	}

	if err := ctx.Validate(req); err != nil {
		validationErr := model.ValidationFailedResponse
		validationErr.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, validationErr)
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

// Me returns information about the authenticated user
func (c *AuthController) Me(ctx echo.Context) error {
	claims, ok := ctx.Get("user").(*service.Claims)
	if !ok {
		return ctx.JSON(http.StatusInternalServerError, model.FailedToGetUserClaimsResponse)
	}

	return ctx.JSON(http.StatusOK, model.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
	})
}
