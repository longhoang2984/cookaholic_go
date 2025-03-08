package domain

import (
	"time"
)

type User struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
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
