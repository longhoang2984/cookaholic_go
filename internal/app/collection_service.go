package app

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
)

type collectionService struct {
	collectionRepo interfaces.CollectionRepository
}

func NewCollectionService(repo interfaces.CollectionRepository) *collectionService {
	return &collectionService{
		collectionRepo: repo,
	}
}

func (s *collectionService) CreateCollection(ctx context.Context, input interfaces.CreateCollectionInput) (*domain.Collection, error) {
	collection := &domain.Collection{
		UserID:      input.UserID,
		Name:        input.Name,
		Description: input.Description,
		Image:       *input.Image,
		BaseModel:   &common.BaseModel{ID: uuid.New()},
	}

	err := s.collectionRepo.Create(ctx, collection)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (s *collectionService) GetCollectionByID(ctx context.Context, id uuid.UUID) (*domain.Collection, error) {

	collection, err := s.collectionRepo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (s *collectionService) GetCollectionByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error) {

	collections, err := s.collectionRepo.GetByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return collections, nil
}

func (s *collectionService) UpdateCollection(ctx context.Context, id uuid.UUID, input interfaces.UpdateCollectionInput) (*domain.Collection, error) {

	collection, err := s.collectionRepo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		collection.Name = input.Name
	}

	if input.Description != "" {
		collection.Description = input.Description
	}

	if input.Image != nil {
		collection.Image = *input.Image
	}

	err = s.collectionRepo.Update(ctx, collection)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (s *collectionService) DeleteCollection(ctx context.Context, id uuid.UUID) error {
	return s.collectionRepo.Delete(ctx, id)
}
