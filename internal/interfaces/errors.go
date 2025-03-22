package interfaces

import "errors"

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrOTPExpired         = errors.New("OTP has expired")
	ErrInvalidOTP         = errors.New("invalid OTP")
	ErrRecipeNotFound     = errors.New("recipe not found")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrRatingNotFound     = errors.New("rating not found")
	ErrUnauthorized       = errors.New("unauthorized")
)

// NotFoundError represents a not found error
type NotFoundError struct {
	message string
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) error {
	return &NotFoundError{message: message}
}

// Error returns the error message
func (e *NotFoundError) Error() string {
	return e.message
}

// UnauthorizedError represents an unauthorized error
type UnauthorizedError struct {
	message string
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) error {
	return &UnauthorizedError{message: message}
}

// Error returns the error message
func (e *UnauthorizedError) Error() string {
	return e.message
}
