package models

import (
	"time"

	"github.com/google/uuid"
)

type Testimonial struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PortfolioID  uuid.UUID `gorm:"type:uuid;not null" json:"portfolio_id"`
	AuthorName   string    `gorm:"type:varchar(255);not null" json:"author_name"`
	AuthorRole   string    `gorm:"type:varchar(255)" json:"author_role"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	AvatarURL    string    `gorm:"type:varchar(255)" json:"avatar_url"`
	Order        int       `gorm:"type:integer;default:0" json:"order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}