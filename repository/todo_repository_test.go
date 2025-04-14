package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/repository"
)

func TestTodoRepository(t *testing.T) {
	userRepo := repository.NewUserRepository(testDB)
	todoRepo := repository.NewTodoRepository(testDB)
	ctx := context.Background()

	// Create a user first
	user := &model.User{
		Email:        "todo-test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := userRepo.Create(ctx, user)
	require.NoError(t, err)

	// Test Todo Create
	description := "Test description"
	dueDate := time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
	todo := &model.Todo{
		UserID:      user.ID,
		Title:       "Test Todo",
		Description: &description,
		DueDate:     &dueDate,
		IsCompleted: false,
	}
	err = todoRepo.Create(ctx, todo)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, todo.ID)

	// Test GetByID
	fetchedTodo, err := todoRepo.GetByID(ctx, todo.ID)
	require.NoError(t, err)
	assert.Equal(t, todo.Title, fetchedTodo.Title)
	assert.Equal(t, *todo.Description, *fetchedTodo.Description)
	assert.Equal(t, todo.DueDate.Day(), fetchedTodo.DueDate.Day())
	assert.Equal(t, todo.IsCompleted, fetchedTodo.IsCompleted)

	// Test GetByUserID
	todos, err := todoRepo.GetByUserID(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, todo.ID, todos[0].ID)

	// Test Update
	todo.Title = "Updated Todo"
	newDescription := "Updated description"
	todo.Description = &newDescription
	err = todoRepo.Update(ctx, todo)
	require.NoError(t, err)

	updatedTodo, err := todoRepo.GetByID(ctx, todo.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Todo", updatedTodo.Title)
	assert.Equal(t, "Updated description", *updatedTodo.Description)

	// Test MarkAsCompleted
	err = todoRepo.MarkAsCompleted(ctx, todo.ID)
	require.NoError(t, err)

	completedTodo, err := todoRepo.GetByID(ctx, todo.ID)
	require.NoError(t, err)
	assert.True(t, completedTodo.IsCompleted)

	// Test Delete
	err = todoRepo.Delete(ctx, todo.ID)
	require.NoError(t, err)

	_, err = todoRepo.GetByID(ctx, todo.ID)
	assert.Error(t, err) // Should error as todo is deleted
}
