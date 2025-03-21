package interfaces

import (
	"context"
	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

// RecipeCollectionService defines operations for managing recipe-collection relationships
type RecipeCollectionService interface {
	// SaveRecipeToCollection saves a recipe to a collection
	SaveRecipeToCollection(ctx context.Context, collectionID, recipeID uuid.UUID) error

	// RemoveRecipeFromCollection removes a recipe from a collection
	RemoveRecipeFromCollection(ctx context.Context, collectionID, recipeID uuid.UUID) error

	// GetRecipesByCollectionID retrieves all recipes in a collection with pagination
	GetRecipesByCollectionID(ctx context.Context, collectionID uuid.UUID, limit int, cursor uuid.UUID) ([]domain.Recipe, uuid.UUID, error)

	// GetCollectionsByRecipeID retrieves all collections that contain a recipe
	GetCollectionsByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]domain.Collection, error)

	// IsRecipeInCollection checks if a recipe is in a collection
	IsRecipeInCollection(ctx context.Context, collectionID, recipeID uuid.UUID) (bool, error)
}
