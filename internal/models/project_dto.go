package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	ServiceID   uuid.UUID  `json:"service_id" binding:"required"`
	Title       string     `json:"title" binding:"required,min=3"`
	Summary     string     `json:"summary" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Location    string     `json:"location"`
	StartDate   *time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate     *time.Time `json:"end_date" time_format:"2006-01-02"`
	Status      string     `json:"status" binding:"required,oneof=draft published"`
	IsFeatured  bool       `json:"is_featured"`
	// Karena Jejak Karya butuh visual kuat, kita bisa tangkap Thumbnail dan Galeri
	ThumbnailURL string   `json:"thumbnail_url"` 
	GalleryURLs  []string `json:"gallery_urls"` 
}

type UpdateProjectRequest struct {
	ServiceID   uuid.UUID  `json:"service_id"`
	Title       string     `json:"title" binding:"omitempty,min=3"`
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	StartDate   *time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate     *time.Time `json:"end_date" time_format:"2006-01-02"`
	Status      string     `json:"status" binding:"omitempty,oneof=draft published"`
	IsFeatured  *bool      `json:"is_featured"`
	ThumbnailURL string    `json:"thumbnail_url"`
	GalleryURLs  []string  `json:"gallery_urls"`
}

type ProjectQueryParams struct {
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
	Search     string `form:"search"`
	Status     string `form:"status"`
	ServiceID  string `form:"service_id"` // Penting untuk filter Jejak Karya berdasarkan Layanan
	IsFeatured string `form:"is_featured"`
}