package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status    int       `json:"-" gorm:"column:status;default:1;"`
}

func (c *Category) TableName() string {
	return "categories"
}

// BeforeCreate is a GORM hook that runs before creating a new category
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = now
	c.UpdatedAt = now
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a category
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
