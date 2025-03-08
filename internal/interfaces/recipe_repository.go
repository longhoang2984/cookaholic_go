package interfaces

import (
	"context"
	"cookaholic/internal/domain"
)

type RecipeRepository interface {
	CreateRecipe(ctx context.Context, recipe *domain.Recipe) error
	GetRecipe(ctx context.Context, id uint) (*domain.Recipe, error)
	UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error
	DeleteRecipe(ctx context.Context, id uint) error
	FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uint, limit int) ([]domain.Recipe, uint, error)
}
