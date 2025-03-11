package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
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
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Time        int         `json:"time"` // cooking time in minutes
	CategoryID  uuid.UUID   `json:"category_id"`
	ServingSize int         `json:"serving_size"` // number of people
	Images      StringArray `json:"images"`       // JSON array of image URLs
	Ingredients Ingredients `json:"ingredients"`  // JSON array of ingredients
	Steps       Steps       `json:"steps"`        // JSON array of steps
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Status      int         `json:"status"`
}
