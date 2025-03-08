package db

import (
	"cookaholic/internal/domain"

	"gorm.io/gorm"
)

type MainRepository struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) *MainRepository {
	return &MainRepository{DB: db}
}

func (r *MainRepository) Save(model *domain.Model) error {
	return r.DB.Create(model).Error
}
