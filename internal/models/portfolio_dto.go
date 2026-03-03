package models

type TestimonialRequest struct {
	AuthorName string `json:"author_name" binding:"required"`
	AuthorRole string `json:"author_role"`
	Content    string `json:"content" binding:"required"`
	AvatarURL  string `json:"avatar_url"`
	Order      int    `json:"order"`
}

type CreatePortfolioRequest struct {
	Title        string               `json:"title" binding:"required,min=3"`
	Sector       string               `json:"sector" binding:"required"` // 🌟 HAPUS validasi oneof di baris ini
	ShortStory   string               `json:"short_story" binding:"required"`
	Impact       string               `json:"impact" binding:"required"`
	Location     string               `json:"location" binding:"required"`
	Status       string               `json:"status" binding:"required,oneof=draft published"`
	ThumbnailURL string               `json:"thumbnail_url"`
	GalleryURLs  []string             `json:"gallery_urls"`
	Testimonials []TestimonialRequest `json:"testimonials"`
}

type UpdatePortfolioRequest struct {
	Title        string               `json:"title" binding:"omitempty,min=3"`
	Sector       string               `json:"sector"` // 🌟 HAPUS validasi oneof di baris ini
	ShortStory   string               `json:"short_story"`
	Impact       string               `json:"impact"`
	Location     string               `json:"location"`
	Status       string               `json:"status" binding:"omitempty,oneof=draft published"`
	ThumbnailURL string               `json:"thumbnail_url"`
	GalleryURLs  []string             `json:"gallery_urls"`
	Testimonials []TestimonialRequest `json:"testimonials"`
}

type PortfolioQueryParams struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
	Status string `form:"status"`
	Sector string `form:"sector"` // Sangat penting untuk filter Menu/Submenu
}