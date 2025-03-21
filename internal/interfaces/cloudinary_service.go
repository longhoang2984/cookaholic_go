package interfaces

import (
	"context"
	"cookaholic/internal/common"
	"mime/multipart"
)

// CloudinaryService interface for image upload operations
type CloudinaryService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (*common.Image, error)
	UploadMultipleImages(ctx context.Context, files []*multipart.FileHeader, folder string) ([]*common.Image, error)
}
