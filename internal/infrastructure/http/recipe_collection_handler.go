package http

import (
	"cookaholic/internal/common"
	"cookaholic/internal/infrastructure/db"
	"cookaholic/internal/interfaces"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeCollectionHandler struct {
	recipeCollectionService interfaces.RecipeCollectionService
}

// PaginationResponse is a generic response structure that includes pagination metadata
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	NextCursor string      `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
}

func NewRecipeCollectionHandler(recipeCollectionService interfaces.RecipeCollectionService) *RecipeCollectionHandler {
	return &RecipeCollectionHandler{
		recipeCollectionService: recipeCollectionService,
	}
}

// SaveRecipeToCollection saves a recipe to a collection
func (h *RecipeCollectionHandler) SaveRecipeToCollection(c *gin.Context) {
	collectionIdStr := c.Param("id")
	recipeIdStr := c.Param("recipeId")

	collectionId, err := uuid.Parse(collectionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid collection ID", "InvalidCollectionID"))
		return
	}

	recipeId, err := uuid.Parse(recipeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid recipe ID", "InvalidRecipeID"))
		return
	}

	err = h.recipeCollectionService.SaveRecipeToCollection(c.Request.Context(), collectionId, recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to save recipe to collection", "SaveRecipeToCollectionFailed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe saved to collection successfully"})
}

// RemoveRecipeFromCollection removes a recipe from a collection
func (h *RecipeCollectionHandler) RemoveRecipeFromCollection(c *gin.Context) {
	collectionIdStr := c.Param("id")
	recipeIdStr := c.Param("recipeId")

	collectionId, err := uuid.Parse(collectionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid collection ID", "InvalidCollectionID"))
		return
	}

	recipeId, err := uuid.Parse(recipeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid recipe ID", "InvalidRecipeID"))
		return
	}

	err = h.recipeCollectionService.RemoveRecipeFromCollection(c.Request.Context(), collectionId, recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to remove recipe from collection", "RemoveRecipeFromCollectionFailed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe removed from collection successfully"})
}

// GetRecipesByCollectionID retrieves all recipes in a collection with pagination
func (h *RecipeCollectionHandler) GetRecipesByCollectionID(c *gin.Context) {
	collectionIdStr := c.Param("id")

	collectionId, err := uuid.Parse(collectionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid collection ID", "InvalidCollectionID"))
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	cursorStr := c.DefaultQuery("cursor", "")

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	// Parse cursor
	var cursor uuid.UUID
	if cursorStr != "" {
		cursor, err = uuid.Parse(cursorStr)
		if err != nil {
			cursor = uuid.Nil // Use nil UUID if parsing fails
		}
	} else {
		cursor = uuid.Nil
	}

	// Get recipes with pagination
	recipes, nextCursor, err := h.recipeCollectionService.GetRecipesByCollectionID(c.Request.Context(), collectionId, limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to get recipes by collection ID", "GetRecipesByCollectionIDFailed"))
		return
	}

	// Prepare response with pagination metadata
	response := PaginationResponse{
		Data: recipes,
	}

	// Only set has_more and next_cursor if there are more items
	if nextCursor != uuid.Nil {
		response.HasMore = true
		response.NextCursor = nextCursor.String()
	} else {
		response.HasMore = false
	}

	c.JSON(http.StatusOK, response)
}

// GetCollectionsByRecipeID retrieves all collections that contain a recipe
func (h *RecipeCollectionHandler) GetCollectionsByRecipeID(c *gin.Context) {
	recipeIdStr := c.Param("id")

	recipeId, err := uuid.Parse(recipeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid recipe ID", "InvalidRecipeID"))
		return
	}

	collections, err := h.recipeCollectionService.GetCollectionsByRecipeID(c.Request.Context(), recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to get collections by recipe ID", "GetCollectionsByRecipeIDFailed"))
		return
	}

	c.JSON(http.StatusOK, collections)
}

// IsRecipeInCollection checks if a recipe is in a collection
func (h *RecipeCollectionHandler) IsRecipeInCollection(c *gin.Context) {
	collectionIdStr := c.Param("collectionId")
	recipeIdStr := c.Param("id")

	collectionId, err := uuid.Parse(collectionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid collection ID", "InvalidCollectionID"))
		return
	}

	recipeId, err := uuid.Parse(recipeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewCustomError(err, "Invalid recipe ID", "InvalidRecipeID"))
		return
	}

	isInCollection, err := h.recipeCollectionService.IsRecipeInCollection(c.Request.Context(), collectionId, recipeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewCustomError(err, "Failed to check if recipe is in collection", "IsRecipeInCollectionFailed"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_in_collection": isInCollection})
}

// TestCursorConversion is a debug endpoint to test cursor conversion
func (h *RecipeCollectionHandler) TestCursorConversion(c *gin.Context) {
	// Test the cursor conversion
	testUUID, convertedTime, isValid := db.TestCursorConversion()

	// Format the time for display
	now := time.Now()

	// Prepare response
	response := gin.H{
		"test_uuid":       testUUID.String(),
		"converted_time":  convertedTime.Format(time.RFC3339Nano),
		"current_time":    now.Format(time.RFC3339Nano),
		"is_valid":        isValid,
		"time_difference": now.Sub(convertedTime).String(),
	}

	c.JSON(http.StatusOK, response)
}
