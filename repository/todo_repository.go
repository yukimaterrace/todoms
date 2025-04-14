package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yukimaterrace/todoms/model"
)

// TodoRepository defines the interface for todo data operations
type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Todo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
	MarkAsCompleted(ctx context.Context, id uuid.UUID) error
}

// PostgresTodoRepository implements TodoRepository interface for PostgreSQL
type PostgresTodoRepository struct {
	db *sqlx.DB
}

// NewTodoRepository creates a new PostgresTodoRepository instance
func NewTodoRepository(db *sqlx.DB) TodoRepository {
	return &PostgresTodoRepository{db: db}
}

// Create inserts a new todo into the database
func (r *PostgresTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	if todo.ID == uuid.Nil {
		todo.ID = uuid.New()
	}

	query := `
		INSERT INTO todos (id, user_id, title, description, due_date, is_completed, created_at, updated_at)
		VALUES (:id, :user_id, :title, :description, :due_date, :is_completed, NOW(), NOW())
	`

	_, err := r.db.NamedExecContext(ctx, query, todo)
	return err
}

// GetByID retrieves a todo by its ID
func (r *PostgresTodoRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Todo, error) {
	query := `
		SELECT id, user_id, title, description, due_date, is_completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo model.Todo
	err := r.db.GetContext(ctx, &todo, query, id)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// GetByUserID retrieves all todos for a user
func (r *PostgresTodoRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Todo, error) {
	query := `
		SELECT id, user_id, title, description, due_date, is_completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var todos []model.Todo
	err := r.db.SelectContext(ctx, &todos, query, userID)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

// Update updates an existing todo in the database
func (r *PostgresTodoRepository) Update(ctx context.Context, todo *model.Todo) error {
	query := `
		UPDATE todos
		SET title = :title, description = :description, due_date = :due_date, is_completed = :is_completed
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, todo)
	return err
}

// Delete removes a todo from the database
func (r *PostgresTodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM todos
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MarkAsCompleted sets a todo's is_completed status to true
func (r *PostgresTodoRepository) MarkAsCompleted(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE todos
		SET is_completed = true
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
