package models

import "encoding/json"

// CreatePageRequest membatasi apa saja yang boleh dikirim saat membuat halaman baru
type CreatePageRequest struct {
	Title           string          `json:"title" binding:"required"`
	Slug            string          `json:"slug" binding:"required"`
	TemplateName    string          `json:"template_name" binding:"required"`
	
	// Kita gunakan json.RawMessage agar Gin framework bisa membaca input objek JSON dari frontend
	ContentJSON     json.RawMessage `json:"content_json" binding:"required"` 
	
	MetaTitle       string          `json:"meta_title"`
	MetaDescription string          `json:"meta_description"`
	Status          string          `json:"status" binding:"required,oneof=draft published"`
}

// UpdatePageRequest hampir sama, namun semua field bersifat opsional
type UpdatePageRequest struct {
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	TemplateName    string          `json:"template_name"`
	ContentJSON     json.RawMessage `json:"content_json"`
	MetaTitle       string          `json:"meta_title"`
	MetaDescription string          `json:"meta_description"`
	Status          string          `json:"status" binding:"omitempty,oneof=draft published"`
}