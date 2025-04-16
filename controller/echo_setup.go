package controller

import (
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/yukimaterrace/todoms/handler"
	"github.com/yukimaterrace/todoms/service"
)

// SetupEcho initializes and configures Echo instance with given services
func SetupEcho(userService service.UserService, authService service.AuthenticationService, todoService service.TodoService) *echo.Echo {
	// Initialize Echo
	e := echo.New()
	e.Validator = NewValidator()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Create auth handler
	authHandler := handler.NewAuthHandler(authService)

	// Initialize controllers
	authController := NewAuthController(authService, userService, authHandler)
	todoController := NewTodoController(todoService, authHandler)

	// Register routes
	authController.RegisterRoutes(e)
	todoController.RegisterRoutes(e)

	// Default route
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to todoms API!")
	})

	return e
}
