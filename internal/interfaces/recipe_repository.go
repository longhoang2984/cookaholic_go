package interfaces

import (
	"context"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type RecipeRepository interface {
	CreateRecipe(ctx context.Context, recipe *domain.Recipe) error
	GetRecipe(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
	UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error
	DeleteRecipe(ctx context.Context, id uuid.UUID) error
	FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uuid.UUID, limit int) ([]domain.Recipe, uuid.UUID, error)
}
