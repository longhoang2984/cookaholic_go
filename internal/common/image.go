package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Image struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CloudName string `json:"cloud_name"`
	Extension string `json:"extension"`
}

// Value implements the driver.Valuer interface for Image
func (i Image) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements the sql.Scanner interface for Image
func (i *Image) Scan(value interface{}) error {
	if value == nil {
		*i = Image{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, i)
}
