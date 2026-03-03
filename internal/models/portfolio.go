package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Portfolio struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title             string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug              string         `gorm:"type:varchar(255);unique;not null" json:"slug"`
	Sector            string         `gorm:"type:varchar(100);not null" json:"sector"`
	ShortStory        string         `gorm:"type:text" json:"short_story"`
	Impact            string         `gorm:"type:text" json:"impact"`
	Location          string         `gorm:"type:varchar(255)" json:"location"`
	Status            string         `gorm:"type:varchar(50);default:'draft'" json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi Polimorfik untuk Gambar Utama & Galeri
	Media []MediaAsset `gorm:"polymorphic:Model;polymorphicValue:Portfolio" json:"media"`
	Testimonials []Testimonial `gorm:"foreignKey:PortfolioID;constraint:OnDelete:CASCADE" json:"testimonials"`
}