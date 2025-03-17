package interfaces

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type CategoryService interface {
	Create(ctx context.Context, category CreateCategoryInput) (*domain.Category, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	Update(ctx context.Context, id uuid.UUID, category UpdateCategoryInput) (*domain.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, cursor uuid.UUID, limit int) ([]domain.Category, uuid.UUID, error)
}

type CreateCategoryInput struct {
	Name  string       `json:"name" binding:"required"`
	Image common.Image `json:"image" binding:"required"`
}

type UpdateCategoryInput struct {
	Name  string `json:"name"`
	Image common.Image `json:"image"`
}
