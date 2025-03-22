package http

import (
	"cookaholic/internal/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RecipeRatingHandler handles HTTP requests for recipe ratings
type RecipeRatingHandler struct {
	recipeRatingService interfaces.RecipeRatingService
}

// NewRecipeRatingHandler creates a new RecipeRatingHandler
func NewRecipeRatingHandler(recipeRatingService interfaces.RecipeRatingService) *RecipeRatingHandler {
	return &RecipeRatingHandler{
		recipeRatingService: recipeRatingService,
	}
}

// RateRecipe handles the request to rate a recipe
func (h *RecipeRatingHandler) RateRecipe(c *gin.Context) {
	// Get the authenticated user ID
	uid, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}

	// Parse the recipe ID from the URL parameter
	recipeIDStr := c.Param("id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Bind the request body to the input struct
	var input interfaces.CreateRatingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID and recipe ID from the auth context and URL
	input.UserID = *uid
	input.RecipeID = recipeID

	// Create the rating
	rating, err := h.recipeRatingService.RateRecipe(c.Request.Context(), input)
	if err != nil {
		// Check for specific error types
		switch e := err.(type) {
		case *interfaces.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rating)
}

// UpdateRating handles the request to update a rating
func (h *RecipeRatingHandler) UpdateRating(c *gin.Context) {
	// Get the authenticated user ID
	uid, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}

	// Parse the rating ID from the URL parameter
	ratingIDStr := c.Param("id")
	ratingID, err := uuid.Parse(ratingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// Bind the request body to the input struct
	var input interfaces.UpdateRatingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the rating
	rating, err := h.recipeRatingService.UpdateRating(c.Request.Context(), ratingID, *uid, input)
	if err != nil {
		// Check for specific error types
		switch e := err.(type) {
		case *interfaces.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *interfaces.UnauthorizedError:
			c.JSON(http.StatusForbidden, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rating)
}

// DeleteRating handles the request to delete a rating
func (h *RecipeRatingHandler) DeleteRating(c *gin.Context) {
	// Get the authenticated user ID
	uid, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}

	// Parse the rating ID from the URL parameter
	ratingIDStr := c.Param("id")
	ratingID, err := uuid.Parse(ratingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	// Delete the rating
	if err := h.recipeRatingService.DeleteRating(c.Request.Context(), ratingID, *uid); err != nil {
		// Check for specific error types
		switch e := err.(type) {
		case *interfaces.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *interfaces.UnauthorizedError:
			c.JSON(http.StatusForbidden, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}

// GetRatingsByRecipeID handles the request to get all ratings for a recipe
func (h *RecipeRatingHandler) GetRatingsByRecipeID(c *gin.Context) {
	// Parse the recipe ID from the URL parameter
	recipeIDStr := c.Param("id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Parse pagination parameters
	cursorStr := c.Query("cursor")
	var cursor uuid.UUID
	if cursorStr != "" {
		cursor, err = uuid.Parse(cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
			return
		}
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	// Get the ratings with user information
	ratings, nextCursor, err := h.recipeRatingService.GetRatingsWithUserByRecipeID(c.Request.Context(), recipeID, cursor, limit)
	if err != nil {
		// Check for specific error types
		switch e := err.(type) {
		case *interfaces.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ratings":     ratings,
		"next_cursor": nextCursor,
	})
}

// GetUserRatingForRecipe handles the request to get a user's rating for a recipe
func (h *RecipeRatingHandler) GetUserRatingForRecipe(c *gin.Context) {
	// Get the authenticated user ID
	uid, errResp := AuthorizedPermission(c)
	if errResp != nil {
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}

	// Parse the recipe ID from the URL parameter
	recipeIDStr := c.Param("id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Get the rating
	rating, err := h.recipeRatingService.GetUserRatingForRecipe(c.Request.Context(), *uid, recipeID)
	if err != nil {
		// Check for specific error types
		switch e := err.(type) {
		case *interfaces.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rating)
}
