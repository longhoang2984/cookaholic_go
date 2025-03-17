package http

import (
	"cookaholic/internal/domain"
	"cookaholic/internal/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	collectionService interfaces.CollectionService
}

func NewCollectionHandler(router *gin.Engine, service interfaces.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		collectionService: service,
	}
}

// CreateCollection handles the creation of a new collection
func (h *CollectionHandler) CreateCollection(c *gin.Context) {
	var collection domain.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, err := AuthorizedPermission(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	collection.UserID = *uid

	input := interfaces.CreateCollectionInput{
		UserID:      collection.UserID,
		Name:        collection.Name,
		Description: collection.Description,
		Image:       &collection.Image,
	}

	createdCollection, createErr := h.collectionService.CreateCollection(c.Request.Context(), input)

	if createErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": createErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCollection)
}

// GetCollection retrieves a collection by ID
func (h *CollectionHandler) GetCollection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	collection, err := h.collectionService.GetCollectionByID(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// GetUserCollections retrieves all collections for a user
func (h *CollectionHandler) GetUserCollections(c *gin.Context) {
	uid, authErr := AuthorizedPermission(c)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, authErr)
		return
	}

	collections, err := h.collectionService.GetCollectionByUserID(c.Request.Context(), *uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, collections)
}

// UpdateCollection updates an existing collection
func (h *CollectionHandler) UpdateCollection(c *gin.Context) {
	var collection interfaces.UpdateCollectionInput
	if err := c.ShouldBindJSON(&collection); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	uid, authErr := AuthorizedPermission(c)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, authErr)
		return
	}

	collection.UserID = *uid

	updatedCollection, updateErr := h.collectionService.UpdateCollection(c.Request.Context(), id, collection)

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCollection)
}

// DeleteCollection deletes a collection
func (h *CollectionHandler) DeleteCollection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	uid, authErr := AuthorizedPermission(c)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, authErr)
		return
	}

	collection, err := h.collectionService.GetCollectionByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	if collection.UserID != *uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this collection"})
		return
	}

	if err := h.collectionService.DeleteCollection(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}
