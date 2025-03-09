package db

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecipeRepository struct {
	db *gorm.DB
}

// CreateRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) CreateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	return r.db.WithContext(ctx).Create(recipe).Error
}

// DeleteRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Update("status", 0).Error
}

// GetRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) GetRecipe(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	var recipe domain.Recipe
	if err := r.db.WithContext(ctx).Where("status = ?", 1).First(&recipe, id).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uuid.UUID, limit int) ([]domain.Recipe, uuid.UUID, error) {
	var recipes []domain.Recipe
	var nextCursor uuid.UUID

	query := r.db.WithContext(ctx)

	for key, value := range conditions {
		// Skip if value is empty
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "user_id":
			query = query.Where("user_id = ?", value)
		case "category":
			query = query.Where("category = ?", value)
		case "serving_size":
			query = query.Where("serving_size = ?", value)
		case "ingredients":
			query = query.Where("ingredients = ?", value)
		case "title":
			query = query.Where("title LIKE ?", "%"+value.(string)+"%")
		}
	}

	if err := query.Where("status = ?", 1).Order("created_at DESC").Limit(limit).Find(&recipes).Error; err != nil {
		return nil, uuid.Nil, err
	}

	if len(recipes) > 0 {
		nextCursor = recipes[len(recipes)-1].ID
	}

	return recipes, nextCursor, nil
}

// UpdateRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	// First get the existing recipe to ensure it exists and belongs to the user
	var existingRecipe domain.Recipe
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", recipe.ID, recipe.UserID).First(&existingRecipe).Error; err != nil {
		return err
	}

	// Update the recipe using Save to trigger hooks
	return r.db.WithContext(ctx).Save(recipe).Error
}

func NewRecipeRepository(db *gorm.DB) interfaces.RecipeRepository {
	return &RecipeRepository{db: db}
}
