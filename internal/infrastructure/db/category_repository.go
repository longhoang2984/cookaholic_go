package db

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryEntity struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status    int       `json:"-" gorm:"column:status;default:1;"`
}

func (c *CategoryEntity) TableName() string {
	return "categories"
}

// BeforeCreate is a GORM hook that runs before creating a new category
func (c *CategoryEntity) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = now
	c.UpdatedAt = now
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a category
func (c *CategoryEntity) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

func (c *CategoryEntity) ToCategoryDomain() *domain.Category {
	return &domain.Category{
		ID:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Status:    c.Status,
	}
}

func FromCategoryDomain(category *domain.Category) *CategoryEntity {
	return &CategoryEntity{
		ID:        category.ID,
		Name:      category.Name,
		Image:     category.Image,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
		Status:    category.Status,
	}
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(FromCategoryDomain(category)).Error
}

func (r *categoryRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var category CategoryEntity
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}
	return category.ToCategoryDomain(), nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Save(FromCategoryDomain(category)).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Update("status", 0).Error
}

func (r *categoryRepository) List(ctx context.Context, cursor uuid.UUID, limit int) ([]domain.Category, uuid.UUID, error) {
	var categories []CategoryEntity
	var nextCursor uuid.UUID

	query := r.db.WithContext(ctx)

	if cursor != uuid.Nil {
		query = query.Where("id > ?", cursor)
	}

	query = query.Order("name ASC").Limit(limit)

	if err := query.Find(&categories).Error; err != nil {
		return nil, uuid.Nil, err
	}

	if len(categories) > 0 {
		nextCursor = categories[len(categories)-1].ID
	}

	categoriesDomain := make([]domain.Category, len(categories))
	for i, category := range categories {
		categoriesDomain[i] = *category.ToCategoryDomain()
	}

	return categoriesDomain, nextCursor, nil
}
