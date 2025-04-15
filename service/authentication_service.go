package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yukimaterrace/todoms/config"
	"github.com/yukimaterrace/todoms/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenType represents the type of JWT token
type TokenType string

const (
	// AccessToken is a short-lived token used for accessing protected resources
	AccessToken TokenType = "access"

	// RefreshToken is a long-lived token used for obtaining new access tokens
	RefreshToken TokenType = "refresh"
)

// Custom errors for authentication service
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrJWTTokenCreation   = errors.New("failed to create JWT token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrInvalidTokenType   = errors.New("invalid token type")
)

// AuthenticationService defines the interface for authentication operations
type AuthenticationService interface {
	// Authenticate validates user credentials and returns a token pair if valid
	Authenticate(ctx context.Context, email, password string) (*TokenPair, error)

	// ValidateToken validates a JWT token and returns the claims
	ValidateToken(tokenString string) (*Claims, error)

	// RefreshToken takes a refresh token and returns a new token pair
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

// JWTAuthService implements the AuthenticationService interface using JWT
type JWTAuthService struct {
	userRepo   repository.UserRepository
	authConfig *config.AuthConfig
	logger     *zap.Logger
}

// NewJWTAuthService creates a new JWT authentication service
func NewJWTAuthService(
	userRepo repository.UserRepository,
	authConfig *config.AuthConfig,
	logger *zap.Logger,
) AuthenticationService {
	return &JWTAuthService{
		userRepo:   userRepo,
		authConfig: authConfig,
		logger:     logger,
	}
}

// Authenticate validates user credentials and returns a token pair if valid
func (s *JWTAuthService) Authenticate(ctx context.Context, email, password string) (*TokenPair, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, ErrUserNotFound
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.logger.Warn("invalid credentials attempt",
			zap.String("email", email),
			zap.Error(err))
		return nil, ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, err := s.generateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		s.logger.Error("failed to generate token pair",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("user authenticated successfully",
		zap.String("email", email),
		zap.String("user_id", user.ID.String()))
	return tokenPair, nil
}

// generateTokenPair creates a new access and refresh token pair
func (s *JWTAuthService) generateTokenPair(userID, email string) (*TokenPair, error) {
	// Create access token
	accessToken, err := s.generateToken(userID, email, AccessToken, s.authConfig.AccessTokenExpiry)
	if err != nil {
		s.logger.Error("failed to generate access token",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	// Create refresh token
	refreshToken, err := s.generateToken(userID, email, RefreshToken, s.authConfig.RefreshTokenExpiry)
	if err != nil {
		s.logger.Error("failed to generate refresh token",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	s.logger.Debug("token pair generated",
		zap.String("user_id", userID),
		zap.String("token_type", "pair"))
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken creates a new JWT token
func (s *JWTAuthService) generateToken(userID, email string, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.authConfig.JWTSecret))
	if err != nil {
		s.logger.Error("failed to sign JWT token",
			zap.String("user_id", userID),
			zap.String("token_type", string(tokenType)),
			zap.Error(err))
		return "", ErrJWTTokenCreation
	}

	s.logger.Debug("token generated",
		zap.String("user_id", userID),
		zap.String("token_type", string(tokenType)),
		zap.Duration("expiry", expiry))
	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTAuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Warn("unexpected signing method",
				zap.String("method", fmt.Sprintf("%v", token.Header["alg"])))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.authConfig.JWTSecret), nil
	})

	if err != nil {
		// Check if token is expired
		if errors.Is(err, jwt.ErrTokenExpired) {
			s.logger.Info("token has expired",
				zap.String("user_id", claims.UserID),
				zap.String("token_type", claims.Type))
			return nil, ErrExpiredToken
		}
		s.logger.Warn("invalid token",
			zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		s.logger.Warn("token validation failed",
			zap.String("user_id", claims.UserID),
			zap.String("token_type", claims.Type))
		return nil, ErrInvalidToken
	}

	s.logger.Debug("token validated successfully",
		zap.String("user_id", claims.UserID),
		zap.String("token_type", claims.Type))
	return claims, nil
}

// RefreshToken takes a refresh token and returns a new token pair
func (s *JWTAuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate the refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		s.logger.Error("refresh token validation failed",
			zap.Error(err))
		return nil, err
	}

	// Ensure it's a refresh token
	if claims.Type != string(RefreshToken) {
		s.logger.Warn("invalid token type for refresh",
			zap.String("user_id", claims.UserID),
			zap.String("expected", string(RefreshToken)),
			zap.String("actual", claims.Type))
		return nil, ErrInvalidTokenType
	}

	// Verify the user still exists
	user, err := s.userRepo.GetByEmail(ctx, claims.Email)
	if err != nil {
		s.logger.Error("user not found during token refresh",
			zap.String("email", claims.Email),
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, ErrUserNotFound
	}

	// Generate a new token pair
	tokenPair, err := s.generateTokenPair(user.ID.String(), user.Email)
	if err != nil {
		s.logger.Error("failed to generate new token pair during refresh",
			zap.String("user_id", claims.UserID),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("token refreshed successfully",
		zap.String("user_id", claims.UserID),
		zap.String("email", claims.Email))
	return tokenPair, nil
}
