package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yukimaterrace/todoms/handler"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
)

// TodoController handles todo-related HTTP requests
type TodoController struct {
	todoService service.TodoService
	authHandler *handler.AuthHandler
}

// NewTodoController creates a new TodoController
func NewTodoController(todoService service.TodoService, authHandler *handler.AuthHandler) *TodoController {
	return &TodoController{
		todoService: todoService,
		authHandler: authHandler,
	}
}

// RegisterRoutes registers the todo routes to the given Echo instance
func (c *TodoController) RegisterRoutes(e *echo.Echo) {
	todos := e.Group("/api/todos", c.authHandler.RequireAuth)
	todos.GET("", c.GetTodos)
	todos.GET("/:id", c.GetTodo)
	todos.POST("", c.CreateTodo)
	todos.PUT("/:id", c.UpdateTodo)
	todos.DELETE("/:id", c.DeleteTodo)
}

// getUUIDFromParam is a helper method to extract and validate a UUID from URL parameters
func (c *TodoController) getUUIDFromParam(ctx echo.Context, paramName string) (uuid.UUID, error) {
	idStr := ctx.Param(paramName)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// getUUIDFromParamWithResponse extracts a UUID from URL parameters and handles the error response
// Returns the UUID and true if successful, or uuid.Nil and false if there was an error
func (c *TodoController) getUUIDFromParamWithResponse(ctx echo.Context, paramName string) (uuid.UUID, bool) {
	id, err := c.getUUIDFromParam(ctx, paramName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.InvalidTodoIDFormatResponse)
		return uuid.Nil, false
	}
	return id, true
}

// handleTodoError handles common error patterns for todo operations
// Returns true if an error was handled, false otherwise
func (c *TodoController) handleTodoError(ctx echo.Context, err error) error {
	switch err {
	case service.ErrTodoNotFound:
		return ctx.JSON(http.StatusNotFound, model.TodoNotFoundResponse)
	case service.ErrUnauthorized:
		return ctx.JSON(http.StatusForbidden, model.NoPermissionToAccessTodoResponse)
	default:
		return ctx.JSON(http.StatusInternalServerError, model.FailedToOperateResponse)
	}
}

// GetTodos returns all todos for the authenticated user
func (c *TodoController) GetTodos(ctx echo.Context) error {
	userID, ok := c.authHandler.GetUserIDFromContextWithResponse(ctx)
	if !ok {
		return nil // Response already sent by GetUserIDFromContextWithResponse
	}

	// Get todos from service
	todos, err := c.todoService.GetTodos(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.FailedToOperateResponse)
	}

	// Return response
	return ctx.JSON(http.StatusOK, model.NewTodoListResponse(todos))
}

// GetTodo returns a specific todo for the authenticated user
func (c *TodoController) GetTodo(ctx echo.Context) error {
	userID, ok := c.authHandler.GetUserIDFromContextWithResponse(ctx)
	if !ok {
		return nil // Response already sent by GetUserIDFromContextWithResponse
	}

	// Parse todo ID from URL parameter
	todoID, ok := c.getUUIDFromParamWithResponse(ctx, "id")
	if !ok {
		return nil // Response already sent by getUUIDFromParamWithResponse
	}

	// Get todo from service
	todo, err := c.todoService.GetTodoByID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		return c.handleTodoError(ctx, err)
	}

	// Return response
	return ctx.JSON(http.StatusOK, model.NewTodoResponse(todo))
}

// CreateTodo creates a new todo for the authenticated user
func (c *TodoController) CreateTodo(ctx echo.Context) error {
	userID, ok := c.authHandler.GetUserIDFromContextWithResponse(ctx)
	if !ok {
		return nil // Response already sent by GetUserIDFromContextWithResponse
	}

	// Bind and validate request
	req := new(model.CreateTodoRequest)
	if err := ValidateRequest(ctx, req); err != nil {
		return err // Error response already sent by ValidateRequest
	}

	// Create todo using service
	todo, err := c.todoService.CreateTodo(ctx.Request().Context(), userID, *req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.FailedToOperateResponse)
	}

	// Return response
	return ctx.JSON(http.StatusCreated, model.NewTodoResponse(todo))
}

// UpdateTodo updates a specific todo for the authenticated user
func (c *TodoController) UpdateTodo(ctx echo.Context) error {
	userID, ok := c.authHandler.GetUserIDFromContextWithResponse(ctx)
	if !ok {
		return nil // Response already sent by GetUserIDFromContextWithResponse
	}

	// Parse todo ID from URL parameter
	todoID, ok := c.getUUIDFromParamWithResponse(ctx, "id")
	if !ok {
		return nil // Response already sent by getUUIDFromParamWithResponse
	}

	// Bind and validate request
	req := new(model.UpdateTodoRequest)
	if err := ValidateRequest(ctx, req); err != nil {
		return err // Error response already sent by ValidateRequest
	}

	// Update todo using service
	todo, err := c.todoService.UpdateTodo(ctx.Request().Context(), userID, todoID, *req)
	if err != nil {
		return c.handleTodoError(ctx, err)
	}

	// Return response
	return ctx.JSON(http.StatusOK, model.NewTodoResponse(todo))
}

// DeleteTodo deletes a specific todo for the authenticated user
func (c *TodoController) DeleteTodo(ctx echo.Context) error {
	userID, ok := c.authHandler.GetUserIDFromContextWithResponse(ctx)
	if !ok {
		return nil // Response already sent by GetUserIDFromContextWithResponse
	}

	// Parse todo ID from URL parameter
	todoID, ok := c.getUUIDFromParamWithResponse(ctx, "id")
	if !ok {
		return nil // Response already sent by getUUIDFromParamWithResponse
	}

	// Delete todo using service
	err := c.todoService.DeleteTodo(ctx.Request().Context(), userID, todoID)
	if err != nil {
		return c.handleTodoError(ctx, err)
	}

	// Return success response
	return ctx.NoContent(http.StatusNoContent)
}
