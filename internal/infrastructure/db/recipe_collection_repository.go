package db

import (
	"context"
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecipeCollectionEntity represents the database model for recipe-collection relationships
type RecipeCollectionEntity struct {
	CollectionID uuid.UUID `gorm:"type:char(36);not null;index:idx_collection_recipe,unique"`
	RecipeID     uuid.UUID `gorm:"type:char(36);not null;index:idx_collection_recipe,unique"`
	CreatedAt    time.Time `gorm:"not null;index:idx_created_at"`
}

// TableName specifies the table name for this entity
func (RecipeCollectionEntity) TableName() string {
	return "recipe_collections"
}

// RecipeCollectionRepository is an implementation of the RecipeCollectionRepository interface
type RecipeCollectionRepository struct {
	db *gorm.DB
}

// NewRecipeCollectionRepository creates a new instance of RecipeCollectionRepository
func NewRecipeCollectionRepository(db *gorm.DB) interfaces.RecipeCollectionRepository {
	return &RecipeCollectionRepository{db: db}
}

// SaveRecipeToCollection saves a recipe to a collection
func (r *RecipeCollectionRepository) SaveRecipeToCollection(ctx context.Context, collectionID, recipeID uuid.UUID) error {
	entity := &RecipeCollectionEntity{
		CollectionID: collectionID,
		RecipeID:     recipeID,
		CreatedAt:    time.Now(),
	}

	// Use a transaction to check if the relationship already exists
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if the relationship already exists
		var count int64
		if err := tx.Model(&RecipeCollectionEntity{}).
			Where("collection_id = ? AND recipe_id = ?", collectionID, recipeID).
			Count(&count).Error; err != nil {
			return err
		}

		// If relationship doesn't exist, create it
		if count == 0 {
			return tx.Create(entity).Error
		}

		// Relationship already exists, nothing to do
		return nil
	})
}

// RemoveRecipeFromCollection removes a recipe from a collection
func (r *RecipeCollectionRepository) RemoveRecipeFromCollection(ctx context.Context, collectionID, recipeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("collection_id = ? AND recipe_id = ?", collectionID, recipeID).
		Delete(&RecipeCollectionEntity{}).Error
}

// ConvertTimeToUUID converts a time.Time to a UUID for cursor-based pagination
// This creates a deterministic UUID based on the timestamp
func ConvertTimeToUUID(t time.Time) uuid.UUID {
	// If the time is zero, return a nil UUID
	if t.IsZero() {
		return uuid.Nil
	}

	// Convert timestamp to nanoseconds
	nanos := t.UnixNano()

	// Create a 16-byte array (size of UUID)
	b := make([]byte, 16)

	// Fill the first 8 bytes with the nanosecond timestamp (big-endian)
	binary.BigEndian.PutUint64(b[:8], uint64(nanos))

	// Set the remaining 8 bytes to zeros
	// This ensures consistent behavior and uniqueness based on timestamp
	for i := 8; i < 16; i++ {
		b[i] = 0
	}

	// Set the UUID version (version 4, variant 2)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 2

	// Create a UUID from the bytes
	id, err := uuid.FromBytes(b)
	if err != nil {
		// In case of an error, return a nil UUID
		return uuid.Nil
	}

	return id
}

// ConvertUUIDToTime converts a UUID back to a time.Time
func ConvertUUIDToTime(id uuid.UUID) time.Time {
	// If nil UUID, return zero time
	if id == uuid.Nil {
		return time.Time{}
	}

	// Convert UUID to string and then parse to bytes
	b, err := id.MarshalBinary()
	if err != nil {
		// In case of an error, return a default time
		return time.Time{}
	}

	// Extract the nanosecond timestamp from the first 8 bytes
	nanos := int64(binary.BigEndian.Uint64(b[:8]))

	// Validate nanos to prevent potential overflow or invalid times
	if nanos < 0 || nanos > time.Now().Add(24*time.Hour).UnixNano() {
		return time.Time{}
	}

	// Convert back to time.Time
	return time.Unix(0, nanos)
}

// GetRecipesByCollectionID retrieves recipes in a collection with pagination
func (r *RecipeCollectionRepository) GetRecipesByCollectionID(ctx context.Context, collectionID uuid.UUID, limit int, cursor uuid.UUID) ([]domain.Recipe, uuid.UUID, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	var recipeEntities []RecipeEntity
	var recipeCollections []RecipeCollectionEntity

	// Set up the base query to get recipe collections
	query := r.db.WithContext(ctx).
		Model(&RecipeCollectionEntity{}).
		Where("collection_id = ?", collectionID).
		Order("created_at DESC")

	// Apply cursor-based pagination if cursor is provided
	if cursor != uuid.Nil {
		// Convert the cursor UUID back to a timestamp
		cursorTime := ConvertUUIDToTime(cursor)
		query = query.Where("created_at < ?", cursorTime)
	}

	// Limit the number of results
	query = query.Limit(limit + 1) // +1 to check if there are more results

	// Execute the query to get recipe collections
	if err := query.Find(&recipeCollections).Error; err != nil {
		return nil, uuid.Nil, err
	}

	// Check if there are more results
	var nextCursor uuid.UUID = uuid.Nil
	hasMore := len(recipeCollections) > limit

	// If there are more results, create a cursor from the last item's timestamp
	if hasMore {
		// The last item in our results will be the one to use for the next cursor
		lastItem := recipeCollections[limit]
		nextCursor = ConvertTimeToUUID(lastItem.CreatedAt)

		// Remove the extra item we fetched
		recipeCollections = recipeCollections[:limit]
	}

	// If no results, return an empty slice and nil UUID
	if len(recipeCollections) == 0 {
		return []domain.Recipe{}, uuid.Nil, nil
	}

	// Extract recipe IDs
	recipeIDs := make([]uuid.UUID, len(recipeCollections))
	for i, rc := range recipeCollections {
		recipeIDs[i] = rc.RecipeID
	}

	// Fetch the actual recipes
	if err := r.db.WithContext(ctx).
		Where("id IN ?", recipeIDs).
		Where("status = ?", 1). // Only active recipes
		Find(&recipeEntities).Error; err != nil {
		return nil, uuid.Nil, err
	}

	// Create a map for faster lookup
	recipeMap := make(map[uuid.UUID]RecipeEntity)
	for _, re := range recipeEntities {
		recipeMap[re.ID] = re
	}

	// Maintain the order from recipe_collections (sorted by created_at DESC)
	recipes := make([]domain.Recipe, 0, len(recipeCollections))
	for _, rc := range recipeCollections {
		if re, ok := recipeMap[rc.RecipeID]; ok {
			recipe := re.ToRecipeDomain()
			recipes = append(recipes, *recipe)
		}
	}

	return recipes, nextCursor, nil
}

// GetCollectionsByRecipeID retrieves all collections that contain a recipe
func (r *RecipeCollectionRepository) GetCollectionsByRecipeID(ctx context.Context, recipeID uuid.UUID) ([]domain.Collection, error) {
	var collectionEntities []CollectionEntity

	// Join recipe_collections with collections to get all collections containing the recipe
	err := r.db.WithContext(ctx).
		Table("collections").
		Joins("JOIN recipe_collections ON collections.id = recipe_collections.collection_id").
		Where("recipe_collections.recipe_id = ?", recipeID).
		Find(&collectionEntities).Error

	if err != nil {
		return nil, err
	}

	// Convert entities to domain models
	collections := make([]domain.Collection, len(collectionEntities))
	for i, entity := range collectionEntities {
		collection := entity.ToCollectionDomain()
		collections[i] = *collection
	}

	return collections, nil
}

// IsRecipeInCollection checks if a recipe is in a collection
func (r *RecipeCollectionRepository) IsRecipeInCollection(ctx context.Context, collectionID, recipeID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&RecipeCollectionEntity{}).
		Where("collection_id = ? AND recipe_id = ?", collectionID, recipeID).
		Count(&count).Error

	return count > 0, err
}

// TestCursorConversion is a debug method to test the cursor conversion logic
// It should be removed in production, but helps verify that the conversion works
func TestCursorConversion() (uuid.UUID, time.Time, bool) {
	// Create a test timestamp
	testTime := time.Now()

	// Convert to UUID
	testUUID := ConvertTimeToUUID(testTime)

	// Convert back to time
	convertedTime := ConvertUUIDToTime(testUUID)

	// Calculate the difference (should be minimal)
	diff := testTime.Sub(convertedTime).Nanoseconds()

	// Check if the difference is within acceptable range (1 millisecond)
	isValid := diff >= -1000000 && diff <= 1000000

	return testUUID, convertedTime, isValid
}
