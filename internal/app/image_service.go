package app

import (
	"context"
	"cookaholic/internal/common"
	"cookaholic/internal/interfaces"
	"mime/multipart"
)

// ImageService implements the image service interface
type ImageService struct {
	cloudinaryService interfaces.CloudinaryService
}

// NewImageService creates a new ImageService instance
func NewImageService(cloudinaryService interfaces.CloudinaryService) *ImageService {
	return &ImageService{
		cloudinaryService: cloudinaryService,
	}
}

// UploadFile uploads a single file
func (s *ImageService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*common.Image, error) {
	return s.cloudinaryService.UploadImage(ctx, file, folder)
}

// UploadMultipleFiles uploads multiple files
func (s *ImageService) UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, folder string) ([]*common.Image, error) {
	return s.cloudinaryService.UploadMultipleImages(ctx, files, folder)
}
