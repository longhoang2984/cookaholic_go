package http

import (
	"cookaholic/internal/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecipeHandler struct {
	recipeService interfaces.RecipeService
}

func NewRecipeHandler(router *gin.Engine, recipeService interfaces.RecipeService) *RecipeHandler {
	handler := &RecipeHandler{
		recipeService: recipeService,
	}

	return handler
}

func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	uid, err := AuthorizedPermission(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var input interfaces.CreateRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID from the token
	input.UserID = *uid

	recipe, createErr := h.recipeService.CreateRecipe(c.Request.Context(), input)
	if createErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": createErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID format"})
		return
	}

	recipe, err := h.recipeService.GetRecipe(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID format"})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	var input interfaces.UpdateRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.recipeService.UpdateRecipe(c.Request.Context(), id, uid, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID format"})
		return
	}

	err = h.recipeService.DeleteRecipe(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

func (h *RecipeHandler) FilterRecipes(c *gin.Context) {
	var cursor uuid.UUID
	var err error

	if cursorStr := c.Query("cursor"); cursorStr != "" {
		cursor, err = uuid.Parse(cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor format"})
			return
		}
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}

	// Initialize input with default values
	input := interfaces.FilterRecipesInput{
		Cursor:     cursor,
		Limit:      limit,
		Conditions: make(map[string]interface{}),
	}

	// First try to get conditions from query parameters
	if c.Query("title") != "" {
		input.Conditions["title"] = c.Query("title")
	}

	if c.Query("category") != "" {
		input.Conditions["category"] = c.Query("category")
	}

	if c.Query("serving_size") != "" {
		input.Conditions["serving_size"] = c.Query("serving_size")
	}

	if c.Query("ingredients") != "" {
		input.Conditions["ingredients"] = c.Query("ingredients")
	}

	if c.Query("time") != "" {
		input.Conditions["time"] = c.Query("time")
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		input.Conditions["user_id"] = userID
	}

	recipes, nextCursor, err := h.recipeService.FilterRecipesByCondition(c.Request.Context(), input.Conditions, input.Cursor, input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes, "nextCursor": nextCursor})
}
