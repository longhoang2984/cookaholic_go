package common

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `json:"status" gorm:"column:status;default:1;"`
}
