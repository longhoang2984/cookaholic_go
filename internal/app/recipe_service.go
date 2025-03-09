package app

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"errors"

	"github.com/google/uuid"
)

type recipeService struct {
	recipeRepo interfaces.RecipeRepository
}

func NewRecipeService(recipeRepo interfaces.RecipeRepository) *recipeService {
	return &recipeService{
		recipeRepo: recipeRepo,
	}
}

func (s *recipeService) CreateRecipe(ctx context.Context, input interfaces.CreateRecipeInput) (*domain.Recipe, error) {
	recipe := &domain.Recipe{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		Time:        input.Time,
		Category:    input.Category,
		ServingSize: input.ServingSize,
		Images:      input.Images,
		Ingredients: input.Ingredients,
		Steps:       input.Steps,
	}

	err := s.recipeRepo.CreateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *recipeService) GetRecipe(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	return s.recipeRepo.GetRecipe(ctx, id)
}

func (s *recipeService) UpdateRecipe(ctx context.Context, id uuid.UUID, userID uuid.UUID, input interfaces.UpdateRecipeInput) (*domain.Recipe, error) {
	// First get the existing recipe
	existingRecipe, err := s.recipeRepo.GetRecipe(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the fields that are provided in the input
	if input.Title != "" {
		existingRecipe.Title = input.Title
	}
	if input.Description != "" {
		existingRecipe.Description = input.Description
	}
	if input.Time != 0 {
		existingRecipe.Time = input.Time
	}
	if input.Category != "" {
		existingRecipe.Category = input.Category
	}
	if input.ServingSize != 0 {
		existingRecipe.ServingSize = input.ServingSize
	}
	if input.Images != nil {
		existingRecipe.Images = input.Images
	}
	if input.Ingredients != nil {
		existingRecipe.Ingredients = input.Ingredients
	}
	if input.Steps != nil {
		existingRecipe.Steps = input.Steps
	}

	// Ensure we're using the correct ID and UserID
	existingRecipe.ID = id
	existingRecipe.UserID = userID

	err = s.recipeRepo.UpdateRecipe(ctx, existingRecipe)
	if err != nil {
		return nil, err
	}

	return existingRecipe, nil
}

func (s *recipeService) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	recipe, err := s.GetRecipe(ctx, id)
	if err != nil {
		return err
	}

	if recipe.Status == 0 {
		return errors.New("recipe not found")
	}

	return s.recipeRepo.DeleteRecipe(ctx, id)
}

func (s *recipeService) FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uuid.UUID, limit int) ([]domain.Recipe, uuid.UUID, error) {
	return s.recipeRepo.FilterRecipesByCondition(ctx, conditions, cursor, limit)
}
