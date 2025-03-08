package http

import (
	"cookaholic/internal/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert user ID to uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	var input interfaces.CreateRecipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID from the token
	input.UserID = uid

	recipe, err := h.recipeService.CreateRecipe(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, recipe)
}

func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recipe, err := h.recipeService.GetRecipe(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert user ID to uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	recipe, err := h.recipeService.GetRecipe(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if recipe.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to update this recipe"})
		return
	}

	var input interfaces.UpdateRecipeInput = interfaces.UpdateRecipeInput{
		Title:       recipe.Title,
		Description: recipe.Description,
		Time:        recipe.Time,
		Category:    recipe.Category,
		ServingSize: recipe.ServingSize,
		Images:      recipe.Images,
		Ingredients: recipe.Ingredients,
		Steps:       recipe.Steps,
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateRecipe, err := h.recipeService.UpdateRecipe(c.Request.Context(), uint(id), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updateRecipe)
}

func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.recipeService.DeleteRecipe(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

func (h *RecipeHandler) FilterRecipes(c *gin.Context) {
	cursor, err := strconv.Atoi(c.Query("cursor"))
	if err != nil {
		cursor = 0
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}

	// Initialize input with default values
	input := interfaces.FilterRecipesInput{
		Cursor:     uint(cursor),
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

	if c.Query("user_id") != "" {
		input.Conditions["user_id"] = c.Query("user_id")
	}

	recipes, nextCursor, err := h.recipeService.FilterRecipesByCondition(c.Request.Context(), input.Conditions, input.Cursor, input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes, "nextCursor": nextCursor})
}
