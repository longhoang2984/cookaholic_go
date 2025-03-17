package interfaces

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type RecipeService interface {
	CreateRecipe(ctx context.Context, input CreateRecipeInput) (*domain.Recipe, error)
	GetRecipe(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
	UpdateRecipe(ctx context.Context, id uuid.UUID, userID uuid.UUID, input UpdateRecipeInput) (*domain.Recipe, error)
	DeleteRecipe(ctx context.Context, id uuid.UUID) error
	FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uuid.UUID, limit int) ([]domain.Recipe, uuid.UUID, error)
}

type CreateRecipeInput struct {
	UserID      uuid.UUID           `json:"-"` // "-" means this field won't be included in JSON
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	Time        int                 `json:"time" binding:"required"`
	CategoryID  uuid.UUID           `json:"category_id" binding:"required"`
	ServingSize int                 `json:"serving_size" binding:"required"`
	Images      []common.Image            `json:"images"`
	Ingredients []domain.Ingredient `json:"ingredients" binding:"required"`
	Steps       []domain.Step       `json:"steps" binding:"required"`
}

type UpdateRecipeInput struct {
	UserID      uuid.UUID           `json:"-"` // "-" means this field won't be included in JSON
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Time        int                 `json:"time"`
	CategoryID  uuid.UUID           `json:"category_id"`
	ServingSize int                 `json:"serving_size"`
	Images      []common.Image            `json:"images"`
	Ingredients []domain.Ingredient `json:"ingredients"`
	Steps       []domain.Step       `json:"steps"`
}

type FilterRecipesInput struct {
	Conditions map[string]interface{} `json:"conditions" gorm:"omitempty"`
	Cursor     uuid.UUID              `json:"cursor"`
	Limit      int                    `json:"limit"`
}
