package model

// SignUpRequest represents the request body for user registration
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the response for user data
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
