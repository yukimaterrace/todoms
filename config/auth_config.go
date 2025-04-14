package config

import (
	"time"
)

// AuthConfig holds authentication related configuration
type AuthConfig struct {
	// JWTSecret is the secret key used to sign JWT tokens
	JWTSecret string

	// AccessTokenExpiry is the duration for which an access token is valid
	AccessTokenExpiry time.Duration

	// RefreshTokenExpiry is the duration for which a refresh token is valid
	RefreshTokenExpiry time.Duration
}

// NewAuthConfig creates a new AuthConfig with the provided parameters
func NewAuthConfig(jwtSecret string, accessTokenExpiry, refreshTokenExpiry time.Duration) *AuthConfig {
	return &AuthConfig{
		JWTSecret:          jwtSecret,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}
}

// DefaultAuthConfig returns a default AuthConfig with sensible defaults
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:          "default-secret-key-change-in-production",
		AccessTokenExpiry:  15 * time.Minute,   // Short-lived access tokens
		RefreshTokenExpiry: 7 * 24 * time.Hour, // 7 days for refresh tokens
	}
}
