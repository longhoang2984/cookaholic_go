package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserFollower represents a follower relationship between users
type UserFollower struct {
	CreatedAt   time.Time `json:"created_at"`
	FollowerID  uuid.UUID `json:"follower_id"`  // The user who is following
	FollowingID uuid.UUID `json:"following_id"` // The user who is being followed
}

// BeforeCreate is a GORM hook that runs before creating a new user follower relationship
func (uf *UserFollower) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	uf.CreatedAt = now
	return nil
}

// UserFollowerWithUser contains follower relationship with follower's basic info
type UserFollowerWithUser struct {
	*UserFollower
	Follower *UserBasicInfo `json:"follower"` // The user who is following
}

// UserFollowingWithUser contains follower relationship with following's basic info
type UserFollowingWithUser struct {
	*UserFollower
	Following *UserBasicInfo `json:"following"` // The user who is being followed
}
