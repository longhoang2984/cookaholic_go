package app

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"time"
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

func (s *recipeService) GetRecipe(ctx context.Context, id uint) (*domain.Recipe, error) {
	return s.recipeRepo.GetRecipe(ctx, id)
}

func (s *recipeService) UpdateRecipe(ctx context.Context, id uint, input interfaces.UpdateRecipeInput) (*domain.Recipe, error) {

	recipe, err := s.GetRecipe(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedRecipe := &domain.Recipe{
		ID:          recipe.ID,
		UserID:      recipe.UserID,
		CreatedAt:   recipe.CreatedAt,
		UpdatedAt:   time.Now(),
		Title:       input.Title,
		Description: input.Description,
		Time:        input.Time,
		Category:    input.Category,
		ServingSize: input.ServingSize,
		Images:      input.Images,
		Ingredients: input.Ingredients,
		Steps:       input.Steps,
	}

	updateErr := s.recipeRepo.UpdateRecipe(ctx, updatedRecipe)
	if updateErr != nil {
		return nil, updateErr
	}

	return updatedRecipe, nil
}

func (s *recipeService) DeleteRecipe(ctx context.Context, id uint) error {
	return s.recipeRepo.DeleteRecipe(ctx, id)
}

func (s *recipeService) FilterRecipesByCondition(ctx context.Context, conditions map[string]interface{}, cursor uint, limit int) ([]domain.Recipe, uint, error) {
	return s.recipeRepo.FilterRecipesByCondition(ctx, conditions, cursor, limit)
}
