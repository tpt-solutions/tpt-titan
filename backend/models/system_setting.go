package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting stores a single key/value system configuration entry.
// Values are stored as text and JSON-encoded by the routes layer so that
// booleans, numbers and nested objects survive a round trip.
type SystemSetting struct {
	Key       string    `json:"key" gorm:"type:varchar(128);primaryKey"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate/Update keeps UpdatedAt fresh.
func (s *SystemSetting) BeforeCreate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *SystemSetting) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}
