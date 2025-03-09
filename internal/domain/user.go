package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID  `json:"id" gorm:"type:char(36);primary_key"`
	Username      string     `json:"username" gorm:"unique;not null"`
	Email         string     `json:"email" gorm:"unique;not null"`
	Password      string     `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	FullName      string     `json:"full_name"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	OTP           string     `json:"-" gorm:"default:null"`
	OTPExpiresAt  time.Time  `json:"-" gorm:"default:null"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"-" gorm:"default:null"`
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
