package interfaces

import (
	"context"
	"cookaholic/internal/domain"
)

type RecipeService interface {
	CreateRecipe(ctx context.Context, input CreateRecipeInput) (*domain.Recipe, error)
	GetRecipe(ctx context.Context, id uint) (*domain.Recipe, error)
	UpdateRecipe(ctx context.Context, id uint, userID uint, input UpdateRecipeInput) (*domain.Recipe, error)
	DeleteRecipe(ctx context.Context, id uint) error
	FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uint, limit int) ([]domain.Recipe, uint, error)
}

type CreateRecipeInput struct {
	UserID      uint                `json:"-"` // "-" means this field won't be included in JSON
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	Time        int                 `json:"time" binding:"required"`
	Category    string              `json:"category" binding:"required"`
	ServingSize int                 `json:"serving_size" binding:"required"`
	Images      []string            `json:"images"`
	Ingredients []domain.Ingredient `json:"ingredients" binding:"required"`
	Steps       []domain.Step       `json:"steps" binding:"required"`
}

type UpdateRecipeInput struct {
	UserID      uint                `json:"-"` // "-" means this field won't be included in JSON
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Time        int                 `json:"time"`
	Category    string              `json:"category"`
	ServingSize int                 `json:"serving_size"`
	Images      []string            `json:"images"`
	Ingredients []domain.Ingredient `json:"ingredients"`
	Steps       []domain.Step       `json:"steps"`
}

type FilterRecipesInput struct {
	Conditions map[string]interface{} `json:"conditions" gorm:"omitempty"`
	Cursor     uint                   `json:"cursor"`
	Limit      int                    `json:"limit"`
}
