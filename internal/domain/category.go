package domain

import (
	"cookaholic/internal/common"
)

type Category struct {
	*common.BaseModel
	Name  string       `json:"name"`
	Image common.Image `json:"image"`
}
