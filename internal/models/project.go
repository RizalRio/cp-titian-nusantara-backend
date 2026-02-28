package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ServiceID   uuid.UUID      `gorm:"type:uuid;not null" json:"service_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string         `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Summary     string         `gorm:"type:text" json:"summary"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"type:varchar(255)" json:"location"`
	StartDate   *time.Time     `gorm:"type:date" json:"start_date"` // Pointer agar bisa null jika belum ditentukan
	EndDate     *time.Time     `gorm:"type:date" json:"end_date"`
	Status      string         `gorm:"type:varchar(50);default:'draft'" json:"status"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Service Service      `gorm:"foreignKey:ServiceID" json:"service"`
	Media   []MediaAsset `gorm:"polymorphic:Model;polymorphicValue:Project" json:"media"` // 🌟 Mendukung multiple images (galeri)
}