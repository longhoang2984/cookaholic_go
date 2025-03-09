package app

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"errors"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo interfaces.CategoryRepository
}

func NewCategoryService(categoryRepo interfaces.CategoryRepository) *categoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(ctx context.Context, input interfaces.CreateCategoryInput) (*domain.Category, error) {
	category := &domain.Category{
		Name:  input.Name,
		Image: input.Image,
	}
	err := s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Get(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return s.categoryRepo.Get(ctx, id)
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, input interfaces.UpdateCategoryInput) (*domain.Category, error) {
	category, err := s.categoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if category.Status == 0 {
		return nil, errors.New("category not found")
	}

	category.Name = input.Name
	category.Image = input.Image

	err = s.categoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	category, err := s.categoryRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	if category.Status == 0 {
		return errors.New("category not found")
	}

	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) List(ctx context.Context, cursor uuid.UUID, limit int) ([]domain.Category, uuid.UUID, error) {
	return s.categoryRepo.List(ctx, cursor, limit)
}
