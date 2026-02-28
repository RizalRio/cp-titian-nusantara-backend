package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string         `gorm:"type:varchar(255);not null;unique" json:"name"`
	Slug             string         `gorm:"type:varchar(255);unique;not null" json:"slug"`
	ShortDescription string         `gorm:"type:text" json:"short_description"`
	Description      string         `gorm:"type:text" json:"description"`
	IconName         string         `gorm:"type:varchar(100)" json:"icon_name"`
	IsFlagship       bool           `gorm:"default:false" json:"is_flagship"`
	Status           string         `gorm:"type:varchar(50);default:'draft'" json:"status"`
	
	// 🌟 FITUR BARU HASIL ROMBAK ARSITEKTUR
	ImpactPoints     []string       `gorm:"type:jsonb;serializer:json" json:"impact_points"`
	CTAText          string         `gorm:"type:varchar(100)" json:"cta_text"`
	CTALink          string         `gorm:"type:varchar(255)" json:"cta_link"`

	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Media    []MediaAsset `gorm:"polymorphic:Model;polymorphicValue:Service" json:"media"`
	Projects []Project    `gorm:"foreignKey:ServiceID" json:"projects"` // 🌟 Relasi One-to-Many ke Jejak Karya
}