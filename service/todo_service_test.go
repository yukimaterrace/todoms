// filepath: /Users/y7matsuo/Project/todoms/service/todo_service_test.go
package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
	"go.uber.org/zap"
)

func TestCreateTodo(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name          string
		userID        uuid.UUID
		request       model.CreateTodoRequest
		setupMock     func(*MockTodoRepository)
		expectedError error
		checkTodo     func(*testing.T, *model.Todo)
	}{
		{
			name:   "Success",
			userID: uuid.New(),
			request: model.CreateTodoRequest{
				Title: "Test Todo",
			},
			setupMock: func(m *MockTodoRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(todo *model.Todo) bool {
					return todo.Title == "Test Todo" &&
						todo.UserID != uuid.Nil &&
						!todo.IsCompleted
				})).Return(nil).Run(func(args mock.Arguments) {
					// Set ID when Create is called, simulating DB behavior
					todo := args.Get(1).(*model.Todo)
					todo.ID = uuid.New()
				})
			},
			expectedError: nil,
			checkTodo: func(t *testing.T, todo *model.Todo) {
				assert.NotEqual(t, uuid.Nil, todo.ID)
				assert.Equal(t, "Test Todo", todo.Title)
				assert.False(t, todo.IsCompleted)
			},
		},
		{
			name:   "With Description and DueDate",
			userID: uuid.New(),
			request: func() model.CreateTodoRequest {
				description := "Test Description"
				dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
				return model.CreateTodoRequest{
					Title:       "Test Todo with Details",
					Description: &description,
					DueDate:     &dueDate,
				}
			}(),
			setupMock: func(m *MockTodoRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(todo *model.Todo) bool {
					return todo.Title == "Test Todo with Details" &&
						todo.Description != nil &&
						todo.DueDate != nil
				})).Return(nil).Run(func(args mock.Arguments) {
					todo := args.Get(1).(*model.Todo)
					todo.ID = uuid.New()
				})
			},
			expectedError: nil,
			checkTodo: func(t *testing.T, todo *model.Todo) {
				assert.NotEqual(t, uuid.Nil, todo.ID)
				assert.Equal(t, "Test Todo with Details", todo.Title)
				assert.NotNil(t, todo.Description)
				assert.Equal(t, "Test Description", *todo.Description)
				assert.NotNil(t, todo.DueDate)
			},
		},
		{
			name:   "Repository Error",
			userID: uuid.New(),
			request: model.CreateTodoRequest{
				Title: "Error Todo",
			},
			setupMock: func(m *MockTodoRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(todo *model.Todo) bool {
					return todo.Title == "Error Todo"
				})).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			checkTodo:     nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockTodoRepository)
			tc.setupMock(mockRepo)

			todoService := service.NewTodoService(mockRepo, logger)

			// Execute
			todo, err := todoService.CreateTodo(ctx, tc.userID, tc.request)

			// Verify
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, todo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, todo)
				if tc.checkTodo != nil {
					tc.checkTodo(t, todo)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTodos(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()
	ctx := context.Background()

	// Test cases
	testCases := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(*MockTodoRepository, uuid.UUID)
		expectedError error
		expectedCount int
	}{
		{
			name:   "Success with Multiple Todos",
			userID: uuid.New(),
			setupMock: func(m *MockTodoRepository, userID uuid.UUID) {
				todos := []model.Todo{
					{
						ID:          uuid.New(),
						UserID:      userID,
						Title:       "Todo 1",
						IsCompleted: false,
					},
					{
						ID:          uuid.New(),
						UserID:      userID,
						Title:       "Todo 2",
						IsCompleted: true,
					},
				}
				m.On("GetByUserID", mock.Anything, userID).Return(todos, nil)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name:   "Success with Empty Todos",
			userID: uuid.New(),
			setupMock: func(m *MockTodoRepository, userID uuid.UUID) {
				m.On("GetByUserID", mock.Anything, userID).Return([]model.Todo{}, nil)
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name:   "Repository Error",
			userID: uuid.New(),
			setupMock: func(m *MockTodoRepository, userID uuid.UUID) {
				m.On("GetByUserID", mock.Anything, userID).Return([]model.Todo{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			expectedCount: 0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockTodoRepository)
			tc.setupMock(mockRepo, tc.userID)

			todoService := service.NewTodoService(mockRepo, logger)

			// Execute
			todos, err := todoService.GetTodos(ctx, tc.userID)

			// Verify
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Empty(t, todos)
			} else {
				assert.NoError(t, err)
				assert.Len(t, todos, tc.expectedCount)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTodoByID(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()
	ctx := context.Background()

	// Setup common test data
	userID := uuid.New()
	anotherUserID := uuid.New() // For testing unauthorized access
	todoID := uuid.New()

	// Test cases
	testCases := []struct {
		name          string
		userID        uuid.UUID
		todoID        uuid.UUID
		setupMock     func(*MockTodoRepository, uuid.UUID, uuid.UUID)
		expectedError error
	}{
		{
			name:   "Success",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				todo := &model.Todo{
					ID:          todoID,
					UserID:      userID,
					Title:       "Test Todo",
					IsCompleted: false,
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)
			},
			expectedError: nil,
		},
		{
			name:   "Todo Not Found",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				m.On("GetByID", mock.Anything, todoID).Return(nil, errors.New("todo not found"))
			},
			expectedError: service.ErrTodoNotFound,
		},
		{
			name:   "Unauthorized Access",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				todo := &model.Todo{
					ID:          todoID,
					UserID:      anotherUserID, // Different from the requesting user
					Title:       "Another User's Todo",
					IsCompleted: false,
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)
			},
			expectedError: service.ErrUnauthorized,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockTodoRepository)
			tc.setupMock(mockRepo, tc.userID, tc.todoID)

			todoService := service.NewTodoService(mockRepo, logger)

			// Execute
			todo, err := todoService.GetTodoByID(ctx, tc.userID, tc.todoID)

			// Verify
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, todo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, todo)
				assert.Equal(t, tc.todoID, todo.ID)
				assert.Equal(t, tc.userID, todo.UserID)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()
	ctx := context.Background()

	// Setup common test data
	userID := uuid.New()
	anotherUserID := uuid.New()
	todoID := uuid.New()

	description := "Updated description"
	dueDate := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	updateRequest := model.UpdateTodoRequest{
		Title:       "Updated Title",
		Description: &description,
		DueDate:     &dueDate,
		IsCompleted: true,
	}

	// Test cases
	testCases := []struct {
		name          string
		userID        uuid.UUID
		todoID        uuid.UUID
		request       model.UpdateTodoRequest
		setupMock     func(*MockTodoRepository, uuid.UUID, uuid.UUID)
		expectedError error
		checkTodo     func(*testing.T, *model.Todo)
	}{
		{
			name:    "Success",
			userID:  userID,
			todoID:  todoID,
			request: updateRequest,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				// First get the todo
				todo := &model.Todo{
					ID:     todoID,
					UserID: userID,
					Title:  "Original Title",
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)

				// Then update it
				m.On("Update", mock.Anything, mock.MatchedBy(func(todo *model.Todo) bool {
					return todo.Title == updateRequest.Title &&
						*todo.Description == *updateRequest.Description &&
						todo.IsCompleted == updateRequest.IsCompleted
				})).Return(nil)
			},
			expectedError: nil,
			checkTodo: func(t *testing.T, todo *model.Todo) {
				assert.Equal(t, updateRequest.Title, todo.Title)
				assert.Equal(t, *updateRequest.Description, *todo.Description)
				assert.Equal(t, updateRequest.IsCompleted, todo.IsCompleted)
			},
		},
		{
			name:    "Todo Not Found",
			userID:  userID,
			todoID:  todoID,
			request: updateRequest,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				m.On("GetByID", mock.Anything, todoID).Return(nil, errors.New("todo not found"))
			},
			expectedError: service.ErrTodoNotFound,
			checkTodo:     nil,
		},
		{
			name:    "Unauthorized Access",
			userID:  userID,
			todoID:  todoID,
			request: updateRequest,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				todo := &model.Todo{
					ID:     todoID,
					UserID: anotherUserID, // Different from the requesting user
					Title:  "Another User's Todo",
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)
			},
			expectedError: service.ErrUnauthorized,
			checkTodo:     nil,
		},
		{
			name:    "Repository Update Error",
			userID:  userID,
			todoID:  todoID,
			request: updateRequest,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				// First get the todo
				todo := &model.Todo{
					ID:     todoID,
					UserID: userID,
					Title:  "Original Title",
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)

				// Then fail on update
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			checkTodo:     nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockTodoRepository)
			tc.setupMock(mockRepo, tc.userID, tc.todoID)

			todoService := service.NewTodoService(mockRepo, logger)

			// Execute
			todo, err := todoService.UpdateTodo(ctx, tc.userID, tc.todoID, tc.request)

			// Verify
			if tc.expectedError != nil {
				assert.Error(t, err)
				if err == service.ErrTodoNotFound || err == service.ErrUnauthorized {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				}
				assert.Nil(t, todo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, todo)
				if tc.checkTodo != nil {
					tc.checkTodo(t, todo)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	// Create a test logger
	logger := zap.NewNop()
	ctx := context.Background()

	// Setup common test data
	userID := uuid.New()
	anotherUserID := uuid.New()
	todoID := uuid.New()

	// Test cases
	testCases := []struct {
		name          string
		userID        uuid.UUID
		todoID        uuid.UUID
		setupMock     func(*MockTodoRepository, uuid.UUID, uuid.UUID)
		expectedError error
	}{
		{
			name:   "Success",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				// First get the todo
				todo := &model.Todo{
					ID:     todoID,
					UserID: userID,
					Title:  "Todo to be deleted",
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)

				// Then delete it
				m.On("Delete", mock.Anything, todoID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Todo Not Found",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				m.On("GetByID", mock.Anything, todoID).Return(nil, errors.New("todo not found"))
			},
			expectedError: service.ErrTodoNotFound,
		},
		{
			name:   "Unauthorized Access",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				todo := &model.Todo{
					ID:     todoID,
					UserID: anotherUserID, // Different from the requesting user
					Title:  "Another User's Todo",
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)
			},
			expectedError: service.ErrUnauthorized,
		},
		{
			name:   "Repository Delete Error",
			userID: userID,
			todoID: todoID,
			setupMock: func(m *MockTodoRepository, userID uuid.UUID, todoID uuid.UUID) {
				// First get the todo
				todo := &model.Todo{
					ID:     todoID,
					UserID: userID,
					Title:  "Todo with delete error",
				}
				m.On("GetByID", mock.Anything, todoID).Return(todo, nil)

				// Then fail on delete
				m.On("Delete", mock.Anything, todoID).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockTodoRepository)
			tc.setupMock(mockRepo, tc.userID, tc.todoID)

			todoService := service.NewTodoService(mockRepo, logger)

			// Execute
			err := todoService.DeleteTodo(ctx, tc.userID, tc.todoID)

			// Verify
			if tc.expectedError != nil {
				assert.Error(t, err)
				if err == service.ErrTodoNotFound || err == service.ErrUnauthorized {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.Equal(t, tc.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
