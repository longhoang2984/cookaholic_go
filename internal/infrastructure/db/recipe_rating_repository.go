package db

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecipeRatingEntity is the database model for recipe ratings
type RecipeRatingEntity struct {
	*common.BaseEntity
	RecipeID uuid.UUID `json:"recipe_id" gorm:"type:char(36);not null;index"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Rating   int       `json:"rating" gorm:"not null"`
	Comment  string    `json:"comment"`
}

// TableName returns the table name for the RecipeRatingEntity
func (r *RecipeRatingEntity) TableName() string {
	return "recipe_ratings"
}

// ToRatingDomain converts a RecipeRatingEntity to a domain.RecipeRating
func (r *RecipeRatingEntity) ToRatingDomain() *domain.RecipeRating {
	return &domain.RecipeRating{
		BaseModel: &common.BaseModel{
			ID:        r.ID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Status:    r.Status,
		},
		RecipeID: r.RecipeID,
		UserID:   r.UserID,
		Rating:   r.Rating,
		Comment:  r.Comment,
	}
}

// FromRatingDomain converts a domain.RecipeRating to a RecipeRatingEntity
func FromRatingDomain(rating *domain.RecipeRating) *RecipeRatingEntity {
	// If the rating is nil, return nil
	if rating == nil {
		return nil
	}

	// If the base model is nil, initialize it
	if rating.BaseModel == nil {
		rating.BaseModel = &common.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		}
	}

	return &RecipeRatingEntity{
		BaseEntity: &common.BaseEntity{
			ID:        rating.ID,
			CreatedAt: rating.CreatedAt,
			UpdatedAt: rating.UpdatedAt,
			Status:    rating.Status,
		},
		RecipeID: rating.RecipeID,
		UserID:   rating.UserID,
		Rating:   rating.Rating,
		Comment:  rating.Comment,
	}
}

// RecipeRatingRepository is the repository implementation for recipe ratings
type RecipeRatingRepository struct {
	db *gorm.DB
}

// CreateRating creates a new rating
func (r *RecipeRatingRepository) CreateRating(ctx context.Context, rating *domain.RecipeRating) error {
	entity := FromRatingDomain(rating)
	result := r.db.Create(entity)
	if result.Error != nil {
		return result.Error
	}

	// Update the rating count and average for the recipe
	return r.UpdateRecipeRatingSummary(ctx, rating.RecipeID)
}

// GetRating gets a rating by ID
func (r *RecipeRatingRepository) GetRating(ctx context.Context, id uuid.UUID) (*domain.RecipeRating, error) {
	var entity RecipeRatingEntity
	result := r.db.First(&entity, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity.ToRatingDomain(), nil
}

// UpdateRating updates an existing rating
func (r *RecipeRatingRepository) UpdateRating(ctx context.Context, rating *domain.RecipeRating) error {
	entity := FromRatingDomain(rating)
	result := r.db.Model(&RecipeRatingEntity{}).Where("id = ?", entity.ID).Updates(map[string]interface{}{
		"rating":     entity.Rating,
		"comment":    entity.Comment,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}

	// Update the rating count and average for the recipe
	return r.UpdateRecipeRatingSummary(ctx, rating.RecipeID)
}

// DeleteRating deletes a rating
func (r *RecipeRatingRepository) DeleteRating(ctx context.Context, id uuid.UUID) error {
	// First, get the rating to know which recipe's summary to update
	var entity RecipeRatingEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return err
	}

	recipeID := entity.RecipeID

	// Delete the rating
	if err := r.db.Delete(&RecipeRatingEntity{}, "id = ?", id).Error; err != nil {
		return err
	}

	// Update the rating count and average for the recipe
	return r.UpdateRecipeRatingSummary(ctx, recipeID)
}

// GetRatingsByRecipeID gets all ratings for a recipe
func (r *RecipeRatingRepository) GetRatingsByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRating, uuid.UUID, error) {
	var entities []RecipeRatingEntity
	var query *gorm.DB

	if cursor == uuid.Nil {
		query = r.db.Where("recipe_id = ?", recipeID).Order("created_at DESC").Limit(limit)
	} else {
		var cursorCreatedAt time.Time
		r.db.Model(&RecipeRatingEntity{}).Where("id = ?", cursor).Select("created_at").Scan(&cursorCreatedAt)
		query = r.db.Where("recipe_id = ? AND created_at < ?", recipeID, cursorCreatedAt).Order("created_at DESC").Limit(limit)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, uuid.Nil, err
	}

	ratings := make([]domain.RecipeRating, len(entities))
	for i, entity := range entities {
		rating := entity.ToRatingDomain()
		ratings[i] = *rating
	}

	var nextCursor uuid.UUID
	if len(entities) == limit {
		nextCursor = entities[len(entities)-1].ID
	} else {
		nextCursor = uuid.Nil
	}

	return ratings, nextCursor, nil
}

// GetRatingByUserAndRecipeID gets a rating by user and recipe ID
func (r *RecipeRatingRepository) GetRatingByUserAndRecipeID(ctx context.Context, userID, recipeID uuid.UUID) (*domain.RecipeRating, error) {
	var entity RecipeRatingEntity
	result := r.db.Where("user_id = ? AND recipe_id = ?", userID, recipeID).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity.ToRatingDomain(), nil
}

// UpdateRecipeRatingSummary calculates and updates the rating summary for a recipe
func (r *RecipeRatingRepository) UpdateRecipeRatingSummary(ctx context.Context, recipeID uuid.UUID) error {
	// Calculate the rating count and average
	var count int64
	var avgRating float64

	// Get the count
	if err := r.db.Model(&RecipeRatingEntity{}).Where("recipe_id = ?", recipeID).Count(&count).Error; err != nil {
		return err
	}

	// If count is 0, set average to 0
	if count == 0 {
		avgRating = 0
	} else {
		// Calculate the average
		if err := r.db.Model(&RecipeRatingEntity{}).Where("recipe_id = ?", recipeID).Select("AVG(rating)").Scan(&avgRating).Error; err != nil {
			return err
		}
	}

	// Update the recipe with the new count and average
	return r.db.Model(&RecipeEntity{}).Where("id = ?", recipeID).Updates(map[string]interface{}{
		"rating_count": count,
		"avg_rating":   avgRating,
		"updated_at":   time.Now(),
	}).Error
}

// GetRatingsWithUserByRecipeID gets all ratings with user information for a recipe
func (r *RecipeRatingRepository) GetRatingsWithUserByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRatingWithUser, uuid.UUID, error) {
	var entities []RecipeRatingEntity
	var query *gorm.DB

	if cursor == uuid.Nil {
		query = r.db.Where("recipe_ratings.recipe_id = ?", recipeID).Order("recipe_ratings.created_at DESC").Limit(limit)
	} else {
		var cursorCreatedAt time.Time
		r.db.Model(&RecipeRatingEntity{}).Where("id = ?", cursor).Select("created_at").Scan(&cursorCreatedAt)
		query = r.db.Where("recipe_ratings.recipe_id = ? AND recipe_ratings.created_at < ?", recipeID, cursorCreatedAt).Order("recipe_ratings.created_at DESC").Limit(limit)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, uuid.Nil, err
	}

	// Create slice for results
	ratingsWithUser := make([]domain.RecipeRatingWithUser, len(entities))

	// For each rating, get the user information
	for i, entity := range entities {
		// Convert the rating entity to domain
		rating := entity.ToRatingDomain()

		// Get user information
		var user UserEntity
		if err := r.db.First(&user, "id = ?", entity.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If user not found, create with default values
				ratingsWithUser[i] = domain.RecipeRatingWithUser{
					RecipeRating: rating,
					User: &domain.UserBasicInfo{
						ID: entity.UserID,
					},
				}
				continue
			}
			return nil, uuid.Nil, err
		}

		// Create user basic info
		userInfo := &domain.UserBasicInfo{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Avatar:   user.Avatar,
		}

		// Add to results
		ratingsWithUser[i] = domain.RecipeRatingWithUser{
			RecipeRating: rating,
			User:         userInfo,
		}
	}

	var nextCursor uuid.UUID
	if len(entities) == limit {
		nextCursor = entities[len(entities)-1].ID
	} else {
		nextCursor = uuid.Nil
	}

	return ratingsWithUser, nextCursor, nil
}

// NewRecipeRatingRepository creates a new recipe rating repository
func NewRecipeRatingRepository(db *gorm.DB) interfaces.RecipeRatingRepository {
	// Ensure the table exists
	db.AutoMigrate(&RecipeRatingEntity{})

	// Update existing recipe table with new fields if needed
	db.Exec("ALTER TABLE recipes ADD COLUMN IF NOT EXISTS rating_count INT DEFAULT 0")
	db.Exec("ALTER TABLE recipes ADD COLUMN IF NOT EXISTS avg_rating FLOAT DEFAULT 0")

	return &RecipeRatingRepository{db: db}
}
