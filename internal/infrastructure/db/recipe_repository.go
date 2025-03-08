package db

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"time"

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
func (r *RecipeRepository) DeleteRecipe(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Delete(&domain.Recipe{
		DeletedAt: &now,
	}, id).Error
}

// GetRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) GetRecipe(ctx context.Context, id uint) (*domain.Recipe, error) {
	var recipe domain.Recipe
	if err := r.db.WithContext(ctx).First(&recipe, id).Error; err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *RecipeRepository) FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uint, limit int) ([]domain.Recipe, uint, error) {
	var recipes []domain.Recipe
	var nextCursor uint

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

	if err := query.Order("created_at DESC").Limit(limit).Find(&recipes).Error; err != nil {
		return nil, 0, err
	}

	if len(recipes) > 0 {
		nextCursor = recipes[len(recipes)-1].ID
	}

	return recipes, nextCursor, nil
}

// UpdateRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	return r.db.WithContext(ctx).Model(&domain.Recipe{}).Where("id = ?", recipe.ID).Where("user_id = ?", recipe.UserID).Updates(recipe).Error
}

func NewRecipeRepository(db *gorm.DB) interfaces.RecipeRepository {
	return &RecipeRepository{db: db}
}
