package interfaces

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type CollectionService interface {
	CreateCollection(ctx context.Context, input CreateCollectionInput) (*domain.Collection, error)
	GetCollectionByID(ctx context.Context, id uuid.UUID) (*domain.Collection, error)
	GetCollectionByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error)
	UpdateCollection(ctx context.Context, id uuid.UUID, input UpdateCollectionInput) (*domain.Collection, error)
	DeleteCollection(ctx context.Context, id uuid.UUID) error
}

type CreateCollectionInput struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Image       *common.Image `json:"image"`
	UserID      uuid.UUID     `json:"user_id"`
}

type UpdateCollectionInput struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Image       *common.Image `json:"image"`
	UserID      uuid.UUID     `json:"user_id"`
}
