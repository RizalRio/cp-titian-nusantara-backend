package models

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan tabel 'users' di database
type User struct {
	ID           string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	
	// json:"-" sangat penting agar password hash tidak pernah ikut terkirim ke Frontend!
	PasswordHash string         `gorm:"type:text;not null" json:"-"` 
	
	RoleID       *string        `gorm:"type:uuid" json:"role_id"` // Pointer (*) karena bisa null
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	
	// Fitur Soft Delete bawaan GORM. Data tidak benar-benar hilang saat di-delete.
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` 
}