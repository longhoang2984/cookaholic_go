package domain

import (
	"time"

	"github.com/google/uuid"
)

// RecipeCollection represents a recipe saved to a collection without using a primary key
// It only stores the collection_id, recipe_id, and created_at timestamp
type RecipeCollection struct {
	CollectionID uuid.UUID `json:"collection_id"`
	RecipeID     uuid.UUID `json:"recipe_id"`
	CreatedAt    time.Time `json:"created_at"`
}
