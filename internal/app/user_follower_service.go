package app

import (
	"context"
	"errors"
	"time"

	"cookaholic/internal/domain"
	"cookaholic/internal/infrastructure/db"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
)

// UserFollowerService implements interfaces.UserFollowerService
type UserFollowerService struct {
	userFollowerRepo *db.UserFollowerRepository
	userRepo         interfaces.UserRepository
}

// NewUserFollowerService creates a new UserFollowerService
func NewUserFollowerService(
	userFollowerRepo *db.UserFollowerRepository,
	userRepo interfaces.UserRepository,
) *UserFollowerService {
	return &UserFollowerService{
		userFollowerRepo: userFollowerRepo,
		userRepo:         userRepo,
	}
}

// FollowUser creates a new follower relationship
func (s *UserFollowerService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	// Validate that both users exist
	_, err := s.userRepo.FindByID(ctx, followerID)
	if err != nil {
		return errors.New("follower user not found")
	}

	_, err = s.userRepo.FindByID(ctx, followingID)
	if err != nil {
		return errors.New("following user not found")
	}

	// Cannot follow yourself
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Create the follower relationship
	userFollower := &domain.UserFollower{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.userFollowerRepo.Create(ctx, userFollower)
}

// UnfollowUser removes a follower relationship
func (s *UserFollowerService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	// Check if the follower relationship exists
	isFollowing, err := s.userFollowerRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	if !isFollowing {
		return errors.New("user is not following this account")
	}

	return s.userFollowerRepo.Delete(ctx, followerID, followingID)
}

// IsFollowing checks if a user is following another user
func (s *UserFollowerService) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return s.userFollowerRepo.IsFollowing(ctx, followerID, followingID)
}

// GetFollowers gets all users who follow a specific user using cursor-based pagination
func (s *UserFollowerService) GetFollowers(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*domain.UserBasicInfo, *time.Time, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	return s.userFollowerRepo.GetFollowersWithInfo(ctx, userID, cursor, limit, s.userRepo)
}

// GetFollowing gets all users whom a specific user follows using cursor-based pagination
func (s *UserFollowerService) GetFollowing(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*domain.UserBasicInfo, *time.Time, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	return s.userFollowerRepo.GetFollowingWithInfo(ctx, userID, cursor, limit, s.userRepo)
}

// GetFollowersCount gets the number of followers for a user
func (s *UserFollowerService) GetFollowersCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.userFollowerRepo.GetFollowersCount(ctx, userID)
}

// GetFollowingCount gets the number of users a user is following
func (s *UserFollowerService) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.userFollowerRepo.GetFollowingCount(ctx, userID)
}

// Ensure UserFollowerService implements the UserFollowerService interface
var _ interfaces.UserFollowerService = (*UserFollowerService)(nil)
