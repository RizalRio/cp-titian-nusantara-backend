package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Action      string     `gorm:"type:varchar(50);not null" json:"action"`
	Module      string     `gorm:"type:varchar(100);not null" json:"module"`
	Description string     `gorm:"type:text" json:"description"`
	OldData     string     `gorm:"type:jsonb" json:"old_data"`
	NewData     string     `gorm:"type:jsonb" json:"new_data"`
	IPAddress   string     `gorm:"type:varchar(50)" json:"ip_address"`
	CreatedAt   time.Time  `json:"created_at"`

	// Relasi
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// DTO untuk pencarian
type ActivityLogQueryParams struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1"`
	Search string `form:"search"`
	Module string `form:"module"`
	Action string `form:"action"`
}