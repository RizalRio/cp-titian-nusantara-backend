package models

import (
	"time"

	"github.com/google/uuid"
)

// Model Database
type SiteSetting struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Key       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Type      string    `gorm:"type:varchar(50)" json:"type"` // Cth: text, image, url
	UpdatedAt time.Time `json:"updated_at"`
}