package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `json:"status" gorm:"column:status;default:1;"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	b.CreatedAt = now
	b.UpdatedAt = now
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	b.UpdatedAt = now
	return nil
}
