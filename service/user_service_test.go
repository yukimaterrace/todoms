package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		email         string
		password      string
		setupMock     func(*MockUserRepository)
		expectedError error
	}{
		{
			name:     "Success",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(m *MockUserRepository) {
				// Mock GetByEmail to return nil, indicating no user with this email exists
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("user not found"))
				// Mock Create to return no error
				m.On("Create", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
					return user.Email == "test@example.com" &&
						bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")) == nil
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "Email Already Exists",
			email:    "existing@example.com",
			password: "password123",
			setupMock: func(m *MockUserRepository) {
				// Mock GetByEmail to return a user, indicating a user with this email exists
				existingUser := &model.User{
					ID:           uuid.New(),
					Email:        "existing@example.com",
					PasswordHash: "hashed_password",
				}
				m.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectedError: service.ErrEmailAlreadyExists,
		},
		{
			name:     "Repository Create Error",
			email:    "error@example.com",
			password: "password123",
			setupMock: func(m *MockUserRepository) {
				// Mock GetByEmail to return nil, indicating no user with this email exists
				m.On("GetByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("user not found"))
				// Mock Create to return an error
				m.On("Create", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
					return user.Email == "error@example.com"
				})).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			tc.setupMock(mockRepo)

			userService := service.NewUserService(mockRepo)

			// Execute
			user, err := userService.CreateUser(context.Background(), tc.email, tc.password)

			// Verify
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.email, user.Email)
				assert.NotEmpty(t, user.PasswordHash)
				assert.NotEqual(t, tc.password, user.PasswordHash) // Password should be hashed
			}

			// Verify all mock expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
