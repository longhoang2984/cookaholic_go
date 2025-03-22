package interfaces

import (
	"context"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type RecipeRatingRepository interface {
	// Create a new rating for a recipe
	CreateRating(ctx context.Context, rating *domain.RecipeRating) error

	// Get a specific rating by ID
	GetRating(ctx context.Context, id uuid.UUID) (*domain.RecipeRating, error)

	// Update an existing rating
	UpdateRating(ctx context.Context, rating *domain.RecipeRating) error

	// Delete a rating
	DeleteRating(ctx context.Context, id uuid.UUID) error

	// Get all ratings for a recipe
	GetRatingsByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRating, uuid.UUID, error)

	// Get all ratings with user information for a recipe
	GetRatingsWithUserByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRatingWithUser, uuid.UUID, error)

	// Get a rating by user and recipe ID
	GetRatingByUserAndRecipeID(ctx context.Context, userID, recipeID uuid.UUID) (*domain.RecipeRating, error)

	// Calculate and update the rating summary (count and average) for a recipe
	UpdateRecipeRatingSummary(ctx context.Context, recipeID uuid.UUID) error
}
