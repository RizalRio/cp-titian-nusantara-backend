package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaAsset struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ModelType string    `gorm:"type:varchar(100);not null;index:idx_model" json:"model_type"`
	ModelID   uuid.UUID `gorm:"type:uuid;not null;index:idx_model" json:"model_id"`
	MediaType string    `gorm:"type:varchar(50);not null" json:"media_type"`
	FileURL   string    `gorm:"type:text;not null" json:"file_url"`
	Caption   string    `gorm:"type:text" json:"caption"`
	Order     int       `gorm:"default:0" json:"order"`
	CreatedAt time.Time `json:"created_at"`
}