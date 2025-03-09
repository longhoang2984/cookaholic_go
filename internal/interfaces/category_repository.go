package interfaces

import (
	"context"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, cursor uuid.UUID, limit int) ([]domain.Category, uuid.UUID, error)
}
