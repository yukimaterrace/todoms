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
