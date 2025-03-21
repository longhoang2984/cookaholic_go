package cloudinary

import (
	"context"
	"cookaholic/internal/common"
	"errors"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService handles image uploads to Cloudinary
type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new Cloudinary service instance
func NewCloudinaryService() (*CloudinaryService, error) {
	// Get Cloudinary credentials from environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, errors.New("missing cloudinary configuration")
	}

	// Create a new Cloudinary instance
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{
		cld: cld,
	}, nil
}

// UploadImage uploads a single image file to Cloudinary
func (s *CloudinaryService) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (*common.Image, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, common.ErrCannotSaveFile(err)
	}
	defer src.Close()

	// Upload the file to Cloudinary
	uploadResult, err := s.cld.Upload.Upload(ctx, src, uploader.UploadParams{
		Folder:      "cookaholic",
		AssetFolder: folder,
	})
	if err != nil {
		return nil, common.ErrCannotSaveFile(err)
	}

	// Create and return image object
	image := &common.Image{
		URL:       uploadResult.SecureURL,
		Width:     uploadResult.Width,
		Height:    uploadResult.Height,
		Extension: uploadResult.Format,
	}

	return image, nil
}

// UploadMultipleImages uploads multiple images to Cloudinary
func (s *CloudinaryService) UploadMultipleImages(ctx context.Context, files []*multipart.FileHeader, folder string) ([]*common.Image, error) {
	if len(files) == 0 {
		return nil, errors.New("no files to upload")
	}

	images := make([]*common.Image, 0, len(files))

	// Upload each file
	for _, file := range files {
		image, err := s.UploadImage(ctx, file, folder)
		if err != nil {
			return images, err // Return any images successfully uploaded along with the error
		}
		images = append(images, image)
	}

	return images, nil
}
