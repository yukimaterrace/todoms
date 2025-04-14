package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/repository"
)

func TestUserRepository(t *testing.T) {
	repo := repository.NewUserRepository(testDB)
	ctx := context.Background()

	// Test Create
	user := &model.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)

	// Test GetByID
	fetchedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.Equal(t, user.PasswordHash, fetchedUser.PasswordHash)

	// Test GetByEmail
	fetchedByEmail, err := repo.GetByEmail(ctx, user.Email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, fetchedByEmail.ID)

	// Test Update
	user.Email = "updated@example.com"
	err = repo.Update(ctx, user)
	require.NoError(t, err)

	updatedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated@example.com", updatedUser.Email)

	// Test Delete
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.Error(t, err) // Should error as user is deleted
}
