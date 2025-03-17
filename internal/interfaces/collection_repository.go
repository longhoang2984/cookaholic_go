package interfaces

import (
	"context"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type CollectionRepository interface {
	Create(ctx context.Context, collection *domain.Collection) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Collection, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error)
	Update(ctx context.Context, collection *domain.Collection) error
	Delete(ctx context.Context, id uuid.UUID) error
}
