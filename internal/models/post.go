package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string         `gorm:"type:varchar(255);unique;not null" json:"slug"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	Excerpt     string         `gorm:"type:text" json:"excerpt"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	AuthorID    uuid.UUID      `gorm:"type:uuid;not null" json:"author_id"`
	Status      string         `gorm:"type:varchar(50);default:'draft'" json:"status"`
	PublishedAt *time.Time     `json:"published_at"` // Menggunakan pointer (*) agar bisa bernilai null
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Sembunyikan field ini dari response JSON frontend

	// 🌟 RELASI EAGER LOADING UNTUK FRONTEND
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
	// Asumsi kamu sudah punya model User. Jika belum ada, buat struct dummy User sementara
	Author   User     `gorm:"foreignKey:AuthorID" json:"author"` 
	// GORM otomatis mengelola tabel pivot `post_tags` untuk relasi Many-to-Many ini
	Tags     []Tag    `gorm:"many2many:post_tags;" json:"tags"`
	// Relasi Polymorphic untuk Media Assets (Gambar, Video, dll)
	Media []MediaAsset `gorm:"polymorphic:Model;polymorphicValue:Post" json:"media"`
}