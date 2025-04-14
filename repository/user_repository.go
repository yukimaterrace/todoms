package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yukimaterrace/todoms/model"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PostgresUserRepository implements UserRepository interface for PostgreSQL
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new PostgresUserRepository instance
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *PostgresUserRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES (:id, :email, :password_hash, NOW(), NOW())
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

// GetByID retrieves a user by their ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user in the database
func (r *PostgresUserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = :email, password_hash = :password_hash
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

// Delete removes a user from the database
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
