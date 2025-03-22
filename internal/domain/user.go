package domain

import (
	"cookaholic/internal/common"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	*common.BaseModel
	Username      string       `json:"username"`
	Email         string       `json:"email"`
	Password      string       `json:"-"` // "-" means this field won't be included in JSON
	FullName      string       `json:"full_name"`
	EmailVerified bool         `json:"email_verified"`
	OTP           *string       `json:"-"`
	OTPExpiresAt  *time.Time    `json:"-"`
	Avatar        common.Image `json:"avatar"`
	Bio           string       `json:"bio"`
}

// BeforeCreate is a GORM hook that runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}
