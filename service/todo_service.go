package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/repository"
	"go.uber.org/zap"
)

var (
	// ErrTodoNotFound is returned when a todo with the specified ID is not found
	ErrTodoNotFound = errors.New("todo not found")

	// ErrUnauthorized is returned when a user is not authorized to access a todo
	ErrUnauthorized = errors.New("not authorized to access this todo")
)

// TodoService defines the interface for todo-related business logic
type TodoService interface {
	// CreateTodo creates a new todo for the specified user
	CreateTodo(ctx context.Context, userID uuid.UUID, req model.CreateTodoRequest) (*model.Todo, error)

	// GetTodos retrieves all todos for the specified user
	GetTodos(ctx context.Context, userID uuid.UUID) ([]model.Todo, error)

	// GetTodoByID retrieves a specific todo by ID, ensuring it belongs to the specified user
	GetTodoByID(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (*model.Todo, error)

	// UpdateTodo updates a specific todo, ensuring it belongs to the specified user
	UpdateTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, req model.UpdateTodoRequest) (*model.Todo, error)

	// DeleteTodo deletes a specific todo, ensuring it belongs to the specified user
	DeleteTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) error
}

// DefaultTodoService implements the TodoService interface
type DefaultTodoService struct {
	todoRepo repository.TodoRepository
	logger   *zap.Logger
}

// NewTodoService creates a new DefaultTodoService instance
func NewTodoService(todoRepo repository.TodoRepository, logger *zap.Logger) TodoService {
	return &DefaultTodoService{
		todoRepo: todoRepo,
		logger:   logger,
	}
}

// CreateTodo creates a new todo for the specified user
func (s *DefaultTodoService) CreateTodo(ctx context.Context, userID uuid.UUID, req model.CreateTodoRequest) (*model.Todo, error) {
	todo := &model.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		IsCompleted: false, // New todos are always not completed
	}

	err := s.todoRepo.Create(ctx, todo)
	if err != nil {
		s.logger.Error("failed to create todo",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("todo created successfully",
		zap.String("user_id", userID.String()),
		zap.String("todo_id", todo.ID.String()))
	return todo, nil
}

// GetTodos retrieves all todos for the specified user
func (s *DefaultTodoService) GetTodos(ctx context.Context, userID uuid.UUID) ([]model.Todo, error) {
	todos, err := s.todoRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get todos",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("retrieved todos successfully",
		zap.String("user_id", userID.String()),
		zap.Int("count", len(todos)))
	return todos, nil
}

// GetTodoByID retrieves a specific todo by ID, ensuring it belongs to the specified user
func (s *DefaultTodoService) GetTodoByID(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (*model.Todo, error) {
	todo, err := s.todoRepo.GetByID(ctx, todoID)
	if err != nil {
		s.logger.Error("failed to get todo",
			zap.String("user_id", userID.String()),
			zap.String("todo_id", todoID.String()),
			zap.Error(err))
		return nil, ErrTodoNotFound
	}

	// Check if the todo belongs to the user
	if todo.UserID != userID {
		s.logger.Warn("unauthorized access attempt to todo",
			zap.String("user_id", userID.String()),
			zap.String("todo_id", todoID.String()),
			zap.String("owner_id", todo.UserID.String()))
		return nil, ErrUnauthorized
	}

	s.logger.Info("retrieved todo successfully",
		zap.String("user_id", userID.String()),
		zap.String("todo_id", todoID.String()))
	return todo, nil
}

// UpdateTodo updates a specific todo, ensuring it belongs to the specified user
func (s *DefaultTodoService) UpdateTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, req model.UpdateTodoRequest) (*model.Todo, error) {
	// Check if the todo exists and belongs to the user
	todo, err := s.GetTodoByID(ctx, userID, todoID)
	if err != nil {
		return nil, err
	}

	// Update the todo
	todo.Title = req.Title
	todo.Description = req.Description
	todo.DueDate = req.DueDate
	todo.IsCompleted = req.IsCompleted

	err = s.todoRepo.Update(ctx, todo)
	if err != nil {
		s.logger.Error("failed to update todo",
			zap.String("user_id", userID.String()),
			zap.String("todo_id", todoID.String()),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("todo updated successfully",
		zap.String("user_id", userID.String()),
		zap.String("todo_id", todoID.String()))
	return todo, nil
}

// DeleteTodo deletes a specific todo, ensuring it belongs to the specified user
func (s *DefaultTodoService) DeleteTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) error {
	// Check if the todo exists and belongs to the user
	_, err := s.GetTodoByID(ctx, userID, todoID)
	if err != nil {
		return err
	}

	err = s.todoRepo.Delete(ctx, todoID)
	if err != nil {
		s.logger.Error("failed to delete todo",
			zap.String("user_id", userID.String()),
			zap.String("todo_id", todoID.String()),
			zap.Error(err))
		return err
	}

	s.logger.Info("todo deleted successfully",
		zap.String("user_id", userID.String()),
		zap.String("todo_id", todoID.String()))
	return nil
}
