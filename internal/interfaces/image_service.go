package interfaces

import (
	"context"
	"cookaholic/internal/common"
	"mime/multipart"
)

type ImageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*common.Image, error)
	UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, folder string) ([]*common.Image, error)
}

type UploadImageInput struct {
	Image *common.Image `json:"image"`
}
