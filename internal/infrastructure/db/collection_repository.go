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

type CollectionEntity struct {
	*common.BaseEntity
	UserID      uuid.UUID    `json:"user_id" gorm:"type:char(36);not null"`
	Name        string       `json:"name" gorm:"not null"`
	Description string       `json:"description"`
	Image       common.Image `json:"image" gorm:"serializer:json;type:text"`
}

func (c *CollectionEntity) TableName() string {
	return "collections"
}

func (c *CollectionEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Status = 1
	return nil
}

func (c *CollectionEntity) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	c.UpdatedAt = now
	return nil
}

func (c *CollectionEntity) ToCollectionDomain() *domain.Collection {
	return &domain.Collection{
		BaseModel:   &common.BaseModel{ID: c.ID, UpdatedAt: c.UpdatedAt, Status: c.Status, CreatedAt: c.CreatedAt},
		UserID:      c.UserID,
		Name:        c.Name,
		Description: c.Description,
		Image:       c.Image,
	}
}

func FromCollectionDomain(collection *domain.Collection) *CollectionEntity {
	return &CollectionEntity{
		BaseEntity:  &common.BaseEntity{ID: collection.ID, UpdatedAt: collection.UpdatedAt, Status: collection.Status},
		UserID:      collection.UserID,
		Name:        collection.Name,
		Description: collection.Description,
		Image:       collection.Image,
	}
}

type CollectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) interfaces.CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) Create(ctx context.Context, collection *domain.Collection) error {
	return r.db.WithContext(ctx).Create(FromCollectionDomain(collection)).Error
}

func (r *CollectionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Collection, error) {
	var collection CollectionEntity
	if err := r.db.WithContext(ctx).First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if collection.Status == 0 {
		return nil, errors.New("collection not found")
	}

	return collection.ToCollectionDomain(), nil
}

func (r *CollectionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error) {
	var collections []CollectionEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Where("status = ?", 1).Find(&collections).Error; err != nil {
		return nil, err
	}

	collectionsDomain := make([]domain.Collection, len(collections))
	for i, collection := range collections {
		collectionsDomain[i] = *collection.ToCollectionDomain()
	}

	return collectionsDomain, nil
}

func (r *CollectionRepository) Update(ctx context.Context, collection *domain.Collection) error {
	existingCollection, err := r.GetByID(ctx, collection.ID)
	if err != nil {
		return err
	}

	if existingCollection.Status == 0 {
		return errors.New("collection not found")
	}

	updatedCollection := FromCollectionDomain(collection)
	existingCollection.Name = updatedCollection.Name
	existingCollection.Description = updatedCollection.Description
	existingCollection.Image = updatedCollection.Image

	updateError := r.db.WithContext(ctx).Where("id = ?", existingCollection.ID).Save(&existingCollection).Error
	if updateError != nil {
		return updateError
	}
	return nil
}

func (r *CollectionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var collection CollectionEntity
	if err := r.db.WithContext(ctx).First(&collection, "id = ?", id).Error; err != nil {
		return err
	}

	if collection.Status == 0 {
		return errors.New("collection not found")
	}

	collection.Status = 0
	return r.db.WithContext(ctx).Save(&collection).Error
}
