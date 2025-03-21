package interfaces

import (
	"context"
	"cookaholic/internal/common"
)

type ImageRepository interface {
	Upload(ctx context.Context, image *common.Image) error
}
