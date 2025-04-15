package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/service"
)

// MockAuthenticationService is a mock of AuthenticationService interface
type MockAuthenticationService struct {
	mock.Mock
}

// Authenticate mocks the Authenticate method
func (m *MockAuthenticationService) Authenticate(ctx context.Context, email, password string) (*service.TokenPair, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenPair), args.Error(1)
}

// ValidateToken mocks the ValidateToken method
func (m *MockAuthenticationService) ValidateToken(tokenString string) (*service.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Claims), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockAuthenticationService) RefreshToken(ctx context.Context, refreshToken string) (*service.TokenPair, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.TokenPair), args.Error(1)
}

func TestRequireAuth(t *testing.T) {
	// Test cases
	tests := []struct {
		name               string
		setupHeader        func(req *http.Request)
		setupMock          func(mockService *MockAuthenticationService)
		expectedStatusCode int
		expectedError      *model.ErrorResponse
		checkContext       bool // True if we want to verify the claims were added to the context
	}{
		{
			name: "Valid Token",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
			},
			setupMock: func(mockService *MockAuthenticationService) {
				claims := &service.Claims{
					UserID: "user123",
					Email:  "test@example.com",
					Type:   string(service.AccessToken),
				}
				mockService.On("ValidateToken", "valid-token").Return(claims, nil)
			},
			expectedStatusCode: http.StatusOK, // Handler should call next and return its result
			checkContext:       true,          // We should check that claims were added to context
		},
		{
			name:               "Missing Authorization Header",
			setupHeader:        func(req *http.Request) {},
			setupMock:          func(mockService *MockAuthenticationService) {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      model.MissingAuthHeaderResponse,
		},
		{
			name: "Invalid Header Format - No Bearer",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "invalid-format")
			},
			setupMock:          func(mockService *MockAuthenticationService) {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      model.InvalidAuthHeaderFormatResponse,
		},
		{
			name: "Invalid Header Format - Empty Token",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer ")
			},
			setupMock:          func(mockService *MockAuthenticationService) {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      model.InvalidAuthHeaderFormatResponse,
		},
		{
			name: "Expired Token",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer expired-token")
			},
			setupMock: func(mockService *MockAuthenticationService) {
				mockService.On("ValidateToken", "expired-token").Return(nil, service.ErrExpiredToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      model.TokenExpiredResponse,
		},
		{
			name: "Invalid Token",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid-token")
			},
			setupMock: func(mockService *MockAuthenticationService) {
				mockService.On("ValidateToken", "invalid-token").Return(nil, errors.New("generic token error"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      model.InvalidTokenResponse,
		},
		{
			name: "Wrong Token Type",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer refresh-token")
			},
			setupMock: func(mockService *MockAuthenticationService) {
				claims := &service.Claims{
					UserID: "user123",
					Email:  "test@example.com",
					Type:   string(service.RefreshToken), // This is a refresh token, not an access token
				}
				mockService.On("ValidateToken", "refresh-token").Return(claims, nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError:      model.InvalidTokenTypeResponse,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Echo
			e := echo.New()

			// Setup mock service
			mockService := new(MockAuthenticationService)
			tc.setupMock(mockService)

			// Create auth handler
			authHandler := NewAuthHandler(mockService)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tc.setupHeader(req)

			// Create response recorder
			rec := httptest.NewRecorder()

			// Create Echo context
			c := e.NewContext(req, rec)

			// Mock handler function that will be called by the middleware if auth succeeds
			nextHandler := func(c echo.Context) error {
				// If we get here, auth succeeded
				// Check if user claims were set in the context
				if tc.checkContext {
					claims, ok := c.Get("user").(*service.Claims)
					assert.True(t, ok)
					assert.Equal(t, "user123", claims.UserID)
					assert.Equal(t, "test@example.com", claims.Email)
				}
				return c.String(http.StatusOK, "Success")
			}

			// Call the middleware
			middleware := authHandler.RequireAuth(nextHandler)
			err := middleware(c)

			// Check error and status
			if tc.expectedStatusCode != http.StatusOK {
				// Authentication failed, expect JSON response in the body
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatusCode, rec.Code)

				// Parse the error response
				var errorResponse model.ErrorResponse
				err = json.Unmarshal(rec.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)

				// Verify the error response matches the expected one
				assert.Equal(t, tc.expectedError.Code, errorResponse.Code)
				assert.Equal(t, tc.expectedError.Message, errorResponse.Message)
			} else {
				// No error should have occurred
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "Success", rec.Body.String())
			}

			// Verify mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestNewAuthHandler(t *testing.T) {
	mockService := new(MockAuthenticationService)
	handler := NewAuthHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.authService)
}
