package db

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Update("status", 0).Error
}

func (r *categoryRepository) List(ctx context.Context, cursor uuid.UUID, limit int) ([]domain.Category, uuid.UUID, error) {
	var categories []domain.Category
	var nextCursor uuid.UUID

	query := r.db.WithContext(ctx)

	query = query.Order("name ASC").Limit(limit)

	if err := query.Find(&categories).Error; err != nil {
		return nil, uuid.Nil, err
	}

	if len(categories) > 0 {
		nextCursor = categories[len(categories)-1].ID
	}

	return categories, nextCursor, nil
}
