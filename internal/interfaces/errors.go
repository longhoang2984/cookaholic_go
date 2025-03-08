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
)
