package db

import (
	"context"
	"errors"
	"time"

	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserFollowerEntity represents the user_followers table in the database
type UserFollowerEntity struct {
	FollowerID  uuid.UUID `gorm:"type:char(36);index:idx_follower"`
	FollowingID uuid.UUID `gorm:"type:char(36);index:idx_following"`
	CreatedAt   time.Time
}

func (UserFollowerEntity) TableName() string {
	return "user_followers"
}

// ToDomain converts UserFollowerEntity to domain.UserFollower
func (e *UserFollowerEntity) ToDomain() *domain.UserFollower {
	return &domain.UserFollower{
		FollowerID:  e.FollowerID,
		FollowingID: e.FollowingID,
		CreatedAt:   e.CreatedAt,
	}
}

// FromDomain converts domain.UserFollower to UserFollowerEntity
func FromUserFollowerDomain(uf *domain.UserFollower) *UserFollowerEntity {
	return &UserFollowerEntity{
		FollowingID: uf.FollowingID,
		FollowerID:  uf.FollowerID,
		CreatedAt:   time.Now(),
	}
}

// UserFollowerRepository provides methods to interact with the user_followers table
type UserFollowerRepository struct {
	db *gorm.DB
}

// NewUserFollowerRepository creates a new UserFollowerRepository
func NewUserFollowerRepository(db *gorm.DB) *UserFollowerRepository {
	return &UserFollowerRepository{db: db}
}

// Create adds a new follower relationship
func (r *UserFollowerRepository) Create(ctx context.Context, userFollower *domain.UserFollower) error {
	entity := FromUserFollowerDomain(userFollower)

	// Check if the relationship already exists
	var existingCount int64
	if err := r.db.Model(&UserFollowerEntity{}).
		Where("follower_id = ? AND following_id = ?",
			entity.FollowerID, entity.FollowingID).
		Count(&existingCount).Error; err != nil {
		return err
	}

	if existingCount > 0 {
		return errors.New("user already follows this account")
	}

	return r.db.Create(&entity).Error
}

// Delete removes a follower relationship (hard delete)
func (r *UserFollowerRepository) Delete(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&UserFollowerEntity{}).Error
}

// IsFollowing checks if a user is following another user
func (r *UserFollowerRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&UserFollowerEntity{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

// GetFollowers gets all users who follow a specific user with cursor-based pagination
func (r *UserFollowerRepository) GetFollowers(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*domain.UserFollower, *time.Time, error) {
	var entities []UserFollowerEntity
	var query = r.db.Where("following_id = ?", userID)

	// Apply cursor if provided
	if cursor != nil {
		query = query.Where("created_at < ?", cursor)
	}

	// Get paginated followers with the most recent first
	if err := query.Order("created_at DESC").
		Limit(limit + 1). // Get one extra to determine if there are more results
		Find(&entities).Error; err != nil {
		return nil, nil, err
	}

	// Determine if there are more results
	var nextCursor *time.Time
	if len(entities) > limit {
		nextCursor = &entities[limit-1].CreatedAt
		entities = entities[:limit] // Remove the extra item
	}

	// Convert to domain objects
	followers := make([]*domain.UserFollower, len(entities))
	for i, entity := range entities {
		followers[i] = entity.ToDomain()
	}

	return followers, nextCursor, nil
}

// GetFollowing gets all users whom a specific user follows with cursor-based pagination
func (r *UserFollowerRepository) GetFollowing(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*domain.UserFollower, *time.Time, error) {
	var entities []UserFollowerEntity
	var query = r.db.Where("follower_id = ?", userID)

	// Apply cursor if provided
	if cursor != nil {
		query = query.Where("created_at < ?", cursor)
	}

	// Get paginated following users with the most recent first
	if err := query.Order("created_at DESC").
		Limit(limit + 1). // Get one extra to determine if there are more results
		Find(&entities).Error; err != nil {
		return nil, nil, err
	}

	// Determine if there are more results
	var nextCursor *time.Time
	if len(entities) > limit {
		nextCursor = &entities[limit-1].CreatedAt
		entities = entities[:limit] // Remove the extra item
	}

	// Convert to domain objects
	following := make([]*domain.UserFollower, len(entities))
	for i, entity := range entities {
		following[i] = entity.ToDomain()
	}

	return following, nextCursor, nil
}

// GetFollowersCount gets the number of followers for a user
func (r *UserFollowerRepository) GetFollowersCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&UserFollowerEntity{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetFollowingCount gets the number of users a user is following
func (r *UserFollowerRepository) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&UserFollowerEntity{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetFollowersWithInfo gets all users who follow a specific user with cursor-based pagination
// and includes their basic user information
func (r *UserFollowerRepository) GetFollowersWithInfo(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int, userRepo interfaces.UserRepository) ([]*domain.UserBasicInfo, *time.Time, error) {
	followers, nextCursor, err := r.GetFollowers(ctx, userID, cursor, limit)
	if err != nil {
		return nil, nil, err
	}

	result := make([]*domain.UserBasicInfo, 0, len(followers))

	// For each follower relationship, get the follower's user info
	for _, follower := range followers {
		// Get follower user info
		user, err := userRepo.FindByID(ctx, follower.FollowerID)
		if err != nil {
			// Skip users that cannot be found (they might have been deleted)
			continue
		}

		// Create UserBasicInfo object
		userInfo := &domain.UserBasicInfo{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Avatar:   user.Avatar,
		}

		// Create combined result
		result = append(result, userInfo)
	}

	return result, nextCursor, nil
}

// GetFollowingWithInfo gets all users whom a specific user follows with cursor-based pagination
// and includes their basic user information
func (r *UserFollowerRepository) GetFollowingWithInfo(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int, userRepo interfaces.UserRepository) ([]*domain.UserBasicInfo, *time.Time, error) {
	following, nextCursor, err := r.GetFollowing(ctx, userID, cursor, limit)
	if err != nil {
		return nil, nil, err
	}

	result := make([]*domain.UserBasicInfo, 0, len(following))

	// For each follower relationship, get the following user's info
	for _, follow := range following {
		// Get following user info
		user, err := userRepo.FindByID(ctx, follow.FollowingID)
		if err != nil {
			// Skip users that cannot be found (they might have been deleted)
			continue
		}

		// Create UserBasicInfo object
		userInfo := &domain.UserBasicInfo{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Avatar:   user.Avatar,
		}

		// Create combined result
		result = append(result, userInfo)
	}

	return result, nextCursor, nil
}
