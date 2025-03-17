package domain

import (
	"cookaholic/internal/common"

	"github.com/google/uuid"
)

type Collection struct {
	*common.BaseModel
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Image       common.Image `json:"image"`
	UserID      uuid.UUID    `json:"user_id"`
}

