package interfaces

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

// UserService defines the interface for user-related operations
type UserService interface {
	Create(ctx context.Context, input CreateUserInput) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, pageSize int) ([]domain.User, error)
	ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error)
	VerifyOTP(ctx context.Context, id uuid.UUID, otp string) error
}

// CreateUserInput defines the input for user creation
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

// UpdateUserInput defines the input for user updates
type UpdateUserInput struct {
	FullName string        `json:"full_name"`
	Password string        `json:"password"`
	Avatar   *common.Image `json:"avatar"`
	Bio      string        `json:"bio"`
}
