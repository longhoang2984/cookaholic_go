package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ingredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type Step struct {
	Order   int    `json:"order"`
	Content string `json:"content"`
}

// Ingredients type for JSON serialization
type Ingredients []Ingredient

// Value implements the driver.Valuer interface for Ingredients
func (i Ingredients) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements the sql.Scanner interface for Ingredients
func (i *Ingredients) Scan(value interface{}) error {
	if value == nil {
		*i = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, i)
}

// Steps type for JSON serialization
type Steps []Step

// Value implements the driver.Valuer interface for Steps
func (s Steps) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for Steps
func (s *Steps) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// StringArray type for JSON serialization of string arrays
type StringArray []string

// Value implements the driver.Valuer interface for StringArray
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for StringArray
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

type Recipe struct {
	ID          uuid.UUID   `json:"id" gorm:"type:char(36);primary_key"`
	UserID      uuid.UUID   `json:"user_id" gorm:"type:char(36);not null"`
	Title       string      `json:"title" gorm:"not null"`
	Description string      `json:"description"`
	Time        int         `json:"time" gorm:"not null"` // cooking time in minutes
	CategoryID  uuid.UUID   `json:"category_id" gorm:"type:char(36);not null"`
	ServingSize int         `json:"serving_size" gorm:"not null"` // number of people
	Images      StringArray `json:"images" gorm:"type:json"`      // JSON array of image URLs
	Ingredients Ingredients `json:"ingredients" gorm:"type:json"` // JSON array of ingredients
	Steps       Steps       `json:"steps" gorm:"type:json"`       // JSON array of steps
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	Status      int         `json:"status" gorm:"column:status;default:1;"`
}

func (r *Recipe) TableName() string {
	return "recipes"
}

// BeforeCreate is a GORM hook that runs before creating a new recipe
func (r *Recipe) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a recipe
func (r *Recipe) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}
