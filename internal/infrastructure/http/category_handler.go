package http

import (
	"cookaholic/internal/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService interfaces.CategoryService
}

func NewCategoryHandler(router *gin.Engine, categoryService interfaces.CategoryService) *CategoryHandler {
	handler := &CategoryHandler{
		categoryService: categoryService,
	}

	return handler
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input interfaces.CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoryService.Create(c.Request.Context(), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.categoryService.Get(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var input interfaces.UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.categoryService.Update(c.Request.Context(), uuid.MustParse(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.categoryService.Delete(c.Request.Context(), uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *CategoryHandler) ListCategories(c *gin.Context) {
	cursor := c.Query("cursor")
	limit := c.Query("limit")

	if cursor == "" {
		cursor = uuid.Nil.String()
	}

	if limit == "" {
		limit = "10"
	}

	cursorUUID, err := uuid.Parse(cursor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	categories, nextCursor, err := h.categoryService.List(c.Request.Context(), cursorUUID, limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories, "nextCursor": nextCursor})

}
