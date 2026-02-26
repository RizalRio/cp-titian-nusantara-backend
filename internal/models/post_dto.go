package models

import "github.com/google/uuid"

// DTO untuk form Create Article dari Admin
type CreatePostRequest struct {
	Title      string      `json:"title" binding:"required,min=5"`
	CategoryID uuid.UUID   `json:"category_id" binding:"required"`
	Excerpt    string      `json:"excerpt"`
	Content    string      `json:"content" binding:"required"`
	Status     string      `json:"status" binding:"required,oneof=draft published"`
	// Menerima array UUID dari Tags yang dipilih di frontend
	TagIDs     []uuid.UUID `json:"tag_ids"` 
	// ThumbnailURL bisa digunakan untuk menyimpan URL gambar utama yang dipilih dari Media Assets
	ThumbnailURL string      `json:"thumbnail_url"`
}

// DTO untuk form Edit Article dari Admin
type UpdatePostRequest struct {
	Title      string      `json:"title" binding:"omitempty,min=5"`
	CategoryID uuid.UUID   `json:"category_id"`
	Excerpt    string      `json:"excerpt"`
	Content    string      `json:"content"`
	Status     string      `json:"status" binding:"omitempty,oneof=draft published"`
	TagIDs     []uuid.UUID `json:"tag_ids"`
	ThumbnailURL string      `json:"thumbnail_url"`
}

// 🌟 DTO Khusus untuk URL Query (Pagination, Filter, & Sort)
// Contoh URL: /api/v1/posts?page=1&limit=10&status=published&category_id=...&search=desa
type PostQueryParams struct {
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
	Search     string `form:"search"`
	Status     string `form:"status"`       // Filter berdasarkan draft/published
	CategoryID string `form:"category_id"`  // Filter berdasarkan kategori
	TagSlug    string `form:"tag_slug"`     // Filter spesifik berdasarkan slug tag
	SortBy     string `form:"sort_by,default=created_at"` // Kolom urutan
	SortOrder  string `form:"sort_order,default=desc"`    // asc atau desc
}