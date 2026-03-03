package models

import (
	"time"

	"github.com/google/uuid"
)

// 🌟 MODELS
type ContactMessage struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Email      string    `gorm:"type:varchar(255);not null" json:"email"`
	Subject    string    `gorm:"type:varchar(255)" json:"subject"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	IsNotified bool      `gorm:"default:false" json:"is_notified"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CollaborationRequest struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationName   string    `gorm:"type:varchar(255);not null" json:"organization_name"`
	ContactPerson      string    `gorm:"type:varchar(255);not null" json:"contact_person"`
	Email              string    `gorm:"type:varchar(255);not null" json:"email"`
	Phone              string    `gorm:"type:varchar(50)" json:"phone"`
	CollaborationType  string    `gorm:"type:varchar(100);not null" json:"collaboration_type"`
	Message            string    `gorm:"type:text;not null" json:"message"`
	ProposalFileURL    string    `gorm:"type:varchar(255)" json:"proposal_file_url"`
	Status             string    `gorm:"type:varchar(50);default:'pending'" json:"status"`
	IsNotified         bool      `gorm:"default:false" json:"is_notified"`
	AssignedTo         *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// 🌟 DTOs (Request Payloads)
type CreateContactRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

type CreateCollaborationRequest struct {
	OrganizationName  string `json:"organization_name" binding:"required"`
	ContactPerson     string `json:"contact_person" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	CollaborationType string `json:"collaboration_type" binding:"required"`
	Message           string `json:"message" binding:"required"`
	ProposalFileURL   string `json:"proposal_file_url"`
}