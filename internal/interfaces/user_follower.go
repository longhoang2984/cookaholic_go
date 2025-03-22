package interfaces

import (
	"context"
	"time"

	"cookaholic/internal/domain"

	"github.com/google/uuid"
)

// UserFollowerService defines the interface for user follower operations
type UserFollowerService interface {
	// FollowUser creates a new follower relationship
	FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error

	// UnfollowUser removes a follower relationship
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error

	// IsFollowing checks if a user is following another user
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)

	// GetFollowers gets all users who follow a specific user using cursor-based pagination
	// The cursor is a timestamp masked as UUID, returns the next cursor for pagination
	// Returns the followers with their basic user information
	GetFollowers(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*domain.UserBasicInfo, *time.Time, error)

	// GetFollowing gets all users whom a specific user follows using cursor-based pagination
	// The cursor is a timestamp masked as UUID, returns the next cursor for pagination
	// Returns the following users with their basic user information
	GetFollowing(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*domain.UserBasicInfo, *time.Time, error)

	// GetFollowersCount gets the number of followers for a user
	GetFollowersCount(ctx context.Context, userID uuid.UUID) (int64, error)

	// GetFollowingCount gets the number of users a user is following
	GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error)
}
