package http

import (
	"cookaholic/internal/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ImageHandler handles HTTP requests for image operations
type ImageHandler struct {
	imageService interfaces.ImageService
}

// NewImageHandler creates a new ImageHandler instance
func NewImageHandler(router *gin.Engine, imageService interfaces.ImageService) *ImageHandler {
	handler := &ImageHandler{
		imageService: imageService,
	}

	return handler
}

// UploadImage handles single image upload
func (h *ImageHandler) UploadImage(c *gin.Context) {
	uid, authErr := AuthorizedPermission(c)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, authErr)
		return
	}

	// Get the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file size (10MB limit)
	if file.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Upload the file
	image, err := h.imageService.UploadFile(c.Request.Context(), file, uid.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": image})
}

// UploadMultipleImages handles multiple image uploads
func (h *ImageHandler) UploadMultipleImages(c *gin.Context) {
	uid, authErr := AuthorizedPermission(c)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, authErr)
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse form"})
		return
	}

	// Get the files from the request
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving files"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	// Check file size limit for each file (10MB per file)
	for _, file := range files {
		if file.Size > 10<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "One or more files are too large"})
			return
		}
	}

	// Upload the files
	images, err := h.imageService.UploadMultipleFiles(c.Request.Context(), files, uid.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}
