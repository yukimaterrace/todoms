package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yukimaterrace/todoms/model"
	"github.com/yukimaterrace/todoms/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrEmailAlreadyExists is returned when trying to create a user with an email that already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// UserService defines the interface for user-related business logic
type UserService interface {
	// CreateUser creates a new user with the given email and password
	CreateUser(ctx context.Context, email, password string) (*model.User, error)
}

// DefaultUserService implements the UserService interface
type DefaultUserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new DefaultUserService instance
func NewUserService(userRepo repository.UserRepository) UserService {
	return &DefaultUserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with the given email and password
func (s *DefaultUserService) CreateUser(ctx context.Context, email, password string) (*model.User, error) {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Generate password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	// Save user to repository
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
