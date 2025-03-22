package app

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

type RecipeRatingService struct {
	recipeRatingRepo interfaces.RecipeRatingRepository
	recipeRepo       interfaces.RecipeRepository
}

// RateRecipe creates a new rating for a recipe
func (s *RecipeRatingService) RateRecipe(ctx context.Context, input interfaces.CreateRatingInput) (*domain.RecipeRating, error) {
	// Check if the recipe exists
	_, err := s.recipeRepo.GetRecipe(ctx, input.RecipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, interfaces.NewNotFoundError("recipe not found")
		}
		return nil, err
	}

	// Check if the user already rated this recipe
	existingRating, err := s.recipeRatingRepo.GetRatingByUserAndRecipeID(ctx, input.UserID, input.RecipeID)
	if err == nil && existingRating != nil {
		// User already rated this recipe, update the existing rating
		existingRating.Rating = input.Rating
		existingRating.Comment = input.Comment
		existingRating.UpdatedAt = time.Now()

		if err := s.recipeRatingRepo.UpdateRating(ctx, existingRating); err != nil {
			return nil, err
		}

		return existingRating, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred
		return nil, err
	}

	// Create a new rating
	rating := &domain.RecipeRating{
		BaseModel: &common.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    1,
		},
		RecipeID: input.RecipeID,
		UserID:   input.UserID,
		Rating:   input.Rating,
		Comment:  input.Comment,
	}

	if err := s.recipeRatingRepo.CreateRating(ctx, rating); err != nil {
		return nil, err
	}

	return rating, nil
}

// UpdateRating updates an existing rating
func (s *RecipeRatingService) UpdateRating(ctx context.Context, id uuid.UUID, userID uuid.UUID, input interfaces.UpdateRatingInput) (*domain.RecipeRating, error) {
	// Get the rating
	rating, err := s.recipeRatingRepo.GetRating(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, interfaces.NewNotFoundError("rating not found")
		}
		return nil, err
	}

	// Check if the user owns the rating
	if rating.UserID != userID {
		return nil, interfaces.NewUnauthorizedError("unauthorized to update this rating")
	}

	// Update the rating
	if input.Rating > 0 {
		rating.Rating = input.Rating
	}
	rating.Comment = input.Comment
	rating.UpdatedAt = time.Now()

	if err := s.recipeRatingRepo.UpdateRating(ctx, rating); err != nil {
		return nil, err
	}

	return rating, nil
}

// DeleteRating deletes a rating
func (s *RecipeRatingService) DeleteRating(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Get the rating
	rating, err := s.recipeRatingRepo.GetRating(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return interfaces.NewNotFoundError("rating not found")
		}
		return err
	}

	// Check if the user owns the rating
	if rating.UserID != userID {
		return interfaces.NewUnauthorizedError("unauthorized to delete this rating")
	}

	// Delete the rating
	return s.recipeRatingRepo.DeleteRating(ctx, id)
}

// GetRatingsByRecipeID gets all ratings for a recipe
func (s *RecipeRatingService) GetRatingsByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRating, uuid.UUID, error) {
	// Check if the recipe exists
	_, err := s.recipeRepo.GetRecipe(ctx, recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, uuid.Nil, interfaces.NewNotFoundError("recipe not found")
		}
		return nil, uuid.Nil, err
	}

	// Get the ratings
	return s.recipeRatingRepo.GetRatingsByRecipeID(ctx, recipeID, cursor, limit)
}

// GetRating gets a rating by ID
func (s *RecipeRatingService) GetRating(ctx context.Context, id uuid.UUID) (*domain.RecipeRating, error) {
	rating, err := s.recipeRatingRepo.GetRating(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, interfaces.NewNotFoundError("rating not found")
		}
		return nil, err
	}
	return rating, nil
}

// GetUserRatingForRecipe gets a user's rating for a recipe
func (s *RecipeRatingService) GetUserRatingForRecipe(ctx context.Context, userID, recipeID uuid.UUID) (*domain.RecipeRating, error) {
	rating, err := s.recipeRatingRepo.GetRatingByUserAndRecipeID(ctx, userID, recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, interfaces.NewNotFoundError("rating not found")
		}
		return nil, err
	}
	return rating, nil
}

// GetRatingsWithUserByRecipeID gets all ratings with user information for a recipe
func (s *RecipeRatingService) GetRatingsWithUserByRecipeID(ctx context.Context, recipeID uuid.UUID, cursor uuid.UUID, limit int) ([]domain.RecipeRatingWithUser, uuid.UUID, error) {
	// Check if the recipe exists
	_, err := s.recipeRepo.GetRecipe(ctx, recipeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, uuid.Nil, interfaces.NewNotFoundError("recipe not found")
		}
		return nil, uuid.Nil, err
	}

	// Get the ratings with user information
	return s.recipeRatingRepo.GetRatingsWithUserByRecipeID(ctx, recipeID, cursor, limit)
}

// NewRecipeRatingService creates a new recipe rating service
func NewRecipeRatingService(recipeRatingRepo interfaces.RecipeRatingRepository, recipeRepo interfaces.RecipeRepository) interfaces.RecipeRatingService {
	return &RecipeRatingService{
		recipeRatingRepo: recipeRatingRepo,
		recipeRepo:       recipeRepo,
	}
}
