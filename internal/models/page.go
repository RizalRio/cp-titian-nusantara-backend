package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Page struct {
	ID              string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title           string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug            string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	TemplateName    string         `gorm:"type:varchar(50);not null" json:"template_name"`
	
	// datatypes.JSON memastikan data tersimpan dengan aman sebagai JSONB di PostgreSQL
	ContentJSON     datatypes.JSON `gorm:"type:jsonb" json:"content_json"` 
	
	MetaTitle       string         `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDescription string         `gorm:"type:varchar(255)" json:"meta_description"`
	Status          string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	
	PublishedAt     *time.Time     `json:"published_at"` // Pointer karena bisa null jika masih draft
	
	// Foreign Keys untuk mencatat siapa admin yang membuat/mengedit
	CreatedBy       string         `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy       *string        `gorm:"type:uuid" json:"updated_by"`
	
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}