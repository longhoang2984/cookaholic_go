package db

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IngredientEntity struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type StepEntity struct {
	Order   int    `json:"order"`
	Content string `json:"content"`
}

// Ingredients type for JSON serialization
type IngredientsEntity []IngredientEntity

// Value implements the driver.Valuer interface for Ingredients
func (i IngredientsEntity) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements the sql.Scanner interface for Ingredients
func (i *IngredientsEntity) Scan(value interface{}) error {
	if value == nil {
		*i = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, i)
}

// Steps type for JSON serialization
type StepsEntity []StepEntity

// Value implements the driver.Valuer interface for Steps
func (s StepsEntity) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for Steps
func (s *StepsEntity) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// StringArray type for JSON serialization of string arrays
type StringArrayEntity []string

// Value implements the driver.Valuer interface for StringArray
func (s StringArrayEntity) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for StringArray
func (s *StringArrayEntity) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

type RecipeEntity struct {
	*common.BaseEntity
	UserID      uuid.UUID         `json:"user_id" gorm:"type:char(36);not null"`
	Title       string            `json:"title" gorm:"not null"`
	Description string            `json:"description"`
	Time        int               `json:"time" gorm:"not null"` // cooking time in minutes
	CategoryID  uuid.UUID         `json:"category_id" gorm:"type:char(36);not null"`
	ServingSize int               `json:"serving_size" gorm:"not null"`            // number of people
	Images      []common.Image    `json:"images" gorm:"serializer:json;type:text"` // JSON array of image URLs
	Ingredients IngredientsEntity `json:"ingredients" gorm:"type:json"`            // JSON array of ingredients
	Steps       StepsEntity       `json:"steps" gorm:"type:json"`                  // JSON array of steps
	RatingCount int               `json:"rating_count" gorm:"default:0"`           // Number of ratings
	AvgRating   float64           `json:"avg_rating" gorm:"default:0"`             // Average rating (0-5)
}

func (r *RecipeEntity) TableName() string {
	return "recipes"
}

func (r *RecipeEntity) ToRecipeDomain() *domain.Recipe {
	ingredients := make([]domain.Ingredient, len(r.Ingredients))
	for i, ingredient := range r.Ingredients {
		ingredients[i] = domain.Ingredient{
			Name:   ingredient.Name,
			Amount: ingredient.Amount,
			Unit:   ingredient.Unit,
		}
	}
	steps := make([]domain.Step, len(r.Steps))
	for i, step := range r.Steps {
		steps[i] = domain.Step{
			Order:   step.Order,
			Content: step.Content,
		}
	}
	var images []common.Image
	if r.Images != nil {
		images = r.Images
	}

	return &domain.Recipe{
		BaseModel: &common.BaseModel{
			ID:        r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Status:    r.Status,
		},
		UserID:      r.UserID,
		Title:       r.Title,
		Description: r.Description,
		Time:        r.Time,
		CategoryID:  r.CategoryID,
		ServingSize: r.ServingSize,
		Images:      images,
		Ingredients: domain.Ingredients(ingredients),
		Steps:       domain.Steps(steps),
		RatingCount: r.RatingCount,
		AvgRating:   r.AvgRating,
	}
}

func FromRecipeDomain(recipe *domain.Recipe) *RecipeEntity {
	// If the recipe is nil, return nil
	if recipe == nil {
		return nil
	}

	ingredients := make([]IngredientEntity, len(recipe.Ingredients))
	for i, ingredient := range recipe.Ingredients {
		ingredients[i] = IngredientEntity{
			Name:   ingredient.Name,
			Amount: ingredient.Amount,
			Unit:   ingredient.Unit,
		}
	}

	steps := make([]StepEntity, len(recipe.Steps))
	for i, step := range recipe.Steps {
		steps[i] = StepEntity{
			Order:   step.Order,
			Content: step.Content,
		}
	}

	var images []common.Image
	if recipe.Images != nil {
		images = recipe.Images
	}

	entity := &RecipeEntity{
		UserID:      recipe.UserID,
		Title:       recipe.Title,
		Description: recipe.Description,
		Time:        recipe.Time,
		CategoryID:  recipe.CategoryID,
		ServingSize: recipe.ServingSize,
		Images:      images,
		Ingredients: ingredients,
		Steps:       steps,
		RatingCount: recipe.RatingCount,
		AvgRating:   recipe.AvgRating,
	}

	if recipe.BaseModel != nil {
		entity.BaseEntity = &common.BaseEntity{
			ID:        recipe.ID,
			CreatedAt: recipe.CreatedAt,
			UpdatedAt: recipe.UpdatedAt,
			Status:    recipe.Status,
		}
	} else {
		entity.BaseEntity = &common.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		}
	}

	return entity
}

type RecipeRepository struct {
	db *gorm.DB
}

// CreateRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) CreateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	return r.db.WithContext(ctx).Create(FromRecipeDomain(recipe)).Error
}

// DeleteRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Update("status", 0).Error
}

// GetRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) GetRecipe(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	var recipe RecipeEntity

	if err := r.db.WithContext(ctx).Where("status = ?", 1).First(&recipe, id).Error; err != nil {
		return nil, err
	}

	return recipe.ToRecipeDomain(), nil
}

func (r *RecipeRepository) FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uuid.UUID, limit int) ([]domain.Recipe, uuid.UUID, error) {
	var recipes []RecipeEntity
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

	recipesDomain := make([]domain.Recipe, len(recipes))
	for i, recipe := range recipes {
		recipesDomain[i] = *recipe.ToRecipeDomain()
	}

	return recipesDomain, nextCursor, nil
}

// UpdateRecipe implements interfaces.RecipeRepository.
func (r *RecipeRepository) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) error {
	// First get the existing recipe to ensure it exists and belongs to the user
	var existingRecipe RecipeEntity
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", recipe.ID, recipe.UserID).First(&existingRecipe).Error; err != nil {
		return err
	}

	if existingRecipe.Status == 0 {
		return errors.New("recipe not found")
	}

	updatedRecipe := FromRecipeDomain(recipe)
	existingRecipe.Title = updatedRecipe.Title
	existingRecipe.Description = updatedRecipe.Description
	existingRecipe.Time = updatedRecipe.Time
	existingRecipe.CategoryID = updatedRecipe.CategoryID
	existingRecipe.ServingSize = updatedRecipe.ServingSize
	existingRecipe.Images = updatedRecipe.Images
	existingRecipe.Ingredients = updatedRecipe.Ingredients
	existingRecipe.Steps = updatedRecipe.Steps

	// Update the recipe using Save to trigger hooks
	return r.db.WithContext(ctx).Save(&existingRecipe).Error
}

func NewRecipeRepository(db *gorm.DB) interfaces.RecipeRepository {
	return &RecipeRepository{db: db}
}
