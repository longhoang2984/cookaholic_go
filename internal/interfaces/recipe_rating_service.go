package interfaces

import (
	"context"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

type RecipeRatingService interface {
	// Rate a recipe
	RateRecipe(ctx context.Context, input CreateRatingInput) (*domain.RecipeRating, error)

	// Update an existing rating
	UpdateRating(ctx context.Context, id uuid.UUID, userID uuid.UUID, input UpdateRatingInput) (*domain.RecipeRating, error)

	// Delete a rating
	DeleteRating(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	// Get all ratings for a recipe
	GetRatingsByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRating, uuid.UUID, error)

	// Get all ratings with user information for a recipe
	GetRatingsWithUserByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRatingWithUser, uuid.UUID, error)

	// Get a specific rating
	GetRating(ctx context.Context, id uuid.UUID) (*domain.RecipeRating, error)

	// Get a user's rating for a recipe
	GetUserRatingForRecipe(ctx context.Context, userID, recipeID uuid.UUID) (*domain.RecipeRating, error)
}

type CreateRatingInput struct {
	UserID   uuid.UUID `json:"user_id"`
	RecipeID uuid.UUID `json:"recipe_id" binding:"required"`
	Rating   int       `json:"rating" binding:"required,min=1,max=5"`
	Comment  string    `json:"comment"`
}

type UpdateRatingInput struct {
	Rating  int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Comment string `json:"comment"`
}
