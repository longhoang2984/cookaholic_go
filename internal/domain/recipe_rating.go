package domain

import (
	"cookaholic/internal/common"

	"github.com/google/uuid"
)

// RecipeRating represents a user's rating and comment for a recipe
type RecipeRating struct {
	*common.BaseModel
	RecipeID uuid.UUID `json:"recipe_id"`
	UserID   uuid.UUID `json:"user_id"`
	Rating   int       `json:"rating"`  // Rating value (1-5)
	Comment  string    `json:"comment"` // Optional comment
}

// RecipeRatingWithUser represents a recipe rating with user information
type RecipeRatingWithUser struct {
	*RecipeRating
	User *UserBasicInfo `json:"user"`
}

// UserBasicInfo contains basic user information for display
type UserBasicInfo struct {
	ID       uuid.UUID    `json:"id"`
	Username string       `json:"username"`
	FullName string       `json:"full_name"`
	Avatar   common.Image `json:"avatar"`
}
