package interfaces

import (
	"context"
	"cookaholic/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]domain.User, error)
}
