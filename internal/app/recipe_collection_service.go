package app

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
)

type recipeCollectionService struct {
	recipeCollectionRepo interfaces.RecipeCollectionRepository
	recipeRepo           interfaces.RecipeRepository
	collectionRepo       interfaces.CollectionRepository
}

// NewRecipeCollectionService creates a new instance of the recipe collection service
func NewRecipeCollectionService(
	recipeCollectionRepo interfaces.RecipeCollectionRepository,
	recipeRepo interfaces.RecipeRepository,
	collectionRepo interfaces.CollectionRepository) interfaces.RecipeCollectionService {
	return &recipeCollectionService{
		recipeCollectionRepo: recipeCollectionRepo,
		recipeRepo:           recipeRepo,
		collectionRepo:       collectionRepo,
	}
}

// SaveRecipeToCollection saves a recipe to a collection
func (s *recipeCollectionService) SaveRecipeToCollection(ctx context.Context, collectionID, recipeID uuid.UUID) error {
	// Verify that the recipe exists
	recipe, err := s.recipeRepo.GetRecipe(ctx, recipeID)
	if err != nil {
		return err
	}
	if recipe == nil {
		return interfaces.ErrRecipeNotFound
	}

	// Verify that the collection exists
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return err
	}
	if collection == nil {
		return interfaces.ErrCollectionNotFound
	}

	// Save the recipe to the collection
	return s.recipeCollectionRepo.SaveRecipeToCollection(ctx, collectionID, recipeID)
}

// RemoveRecipeFromCollection removes a recipe from a collection
func (s *recipeCollectionService) RemoveRecipeFromCollection(ctx context.Context, collectionID, recipeID uuid.UUID) error {
	return s.recipeCollectionRepo.RemoveRecipeFromCollection(ctx, collectionID, recipeID)
}

// GetRecipesByCollectionID retrieves all recipes in a collection with pagination
func (s *recipeCollectionService) GetRecipesByCollectionID(ctx context.Context, collectionID uuid.UUID, limit int, cursor uuid.UUID) ([]domain.Recipe, uuid.UUID, error) {
	// Verify that the collection exists
	collection, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return nil, uuid.Nil, err
	}
	if collection == nil {
		return nil, uuid.Nil, interfaces.ErrCollectionNotFound
	}

	// Set a default limit if not specified or invalid
	if limit <= 0 {
		limit = 10 // Default limit
	} else if limit > 50 {
		limit = 50 // Maximum limit to prevent excessive queries
	}

	// Call the repository with pagination parameters
	return s.recipeCollectionRepo.GetRecipesByCollectionID(ctx, collectionID, limit, cursor)
}

// GetCollectionsByRecipeID retrieves all collections that contain a recipe
func (s *recipeCollectionService) GetCollectionsByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]domain.Collection, error) {
	// Verify that the recipe exists
	recipe, err := s.recipeRepo.GetRecipe(ctx, recipeID)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, interfaces.ErrRecipeNotFound
	}

	return s.recipeCollectionRepo.GetCollectionsByRecipeID(ctx, recipeID)
}

// IsRecipeInCollection checks if a recipe is in a collection
func (s *recipeCollectionService) IsRecipeInCollection(ctx context.Context, collectionID, recipeID uuid.UUID) (bool, error) {
	return s.recipeCollectionRepo.IsRecipeInCollection(ctx, collectionID, recipeID)
}
