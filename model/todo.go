package model

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item in the system
type Todo struct {
	ID          uuid.UUID  `db:"id"`
	UserID      uuid.UUID  `db:"user_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	DueDate     *time.Time `db:"due_date"`
	IsCompleted bool       `db:"is_completed"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// CreateTodoRequest represents the request to create a new todo
type CreateTodoRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
}

// UpdateTodoRequest represents the request to update a todo
type UpdateTodoRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	IsCompleted bool       `json:"isCompleted"`
}

// TodoResponse represents the response for a todo item
type TodoResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	IsCompleted bool       `json:"isCompleted"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// TodoListResponse represents the response for a list of todo items
type TodoListResponse struct {
	Todos []TodoResponse `json:"todos"`
}

// NewTodoResponse creates a new TodoResponse from a Todo model
func NewTodoResponse(todo *Todo) TodoResponse {
	return TodoResponse{
		ID:          todo.ID.String(),
		Title:       todo.Title,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}

// NewTodoListResponse creates a new TodoListResponse from a slice of Todo models
func NewTodoListResponse(todos []Todo) TodoListResponse {
	todoResponses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponses[i] = NewTodoResponse(&todo)
	}
	return TodoListResponse{
		Todos: todoResponses,
	}
}
