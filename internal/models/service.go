package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service merepresentasikan program/layanan utama (Ekosistem Layanan)
type Service struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string         `gorm:"type:varchar(255);not null;unique" json:"name"`
	Slug             string         `gorm:"type:varchar(255);unique;not null" json:"slug"`
	ShortDescription string         `gorm:"type:text" json:"short_description"`
	Description      string         `gorm:"type:text" json:"description"` // Menggunakan description sesuai skema
	IconName         string         `gorm:"type:varchar(100)" json:"icon_name"`
	IsFlagship       bool           `gorm:"default:false" json:"is_flagship"` // Menandakan layanan unggulan
	Status           string         `gorm:"type:varchar(50);default:'draft'" json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi Polimorfik ke Media Assets untuk menangani Thumbnail
	Media []MediaAsset `gorm:"polymorphic:Model;polymorphicValue:Service" json:"media"`
}