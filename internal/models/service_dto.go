package models

// CreateServiceRequest memvalidasi payload saat Admin membuat layanan baru
type CreateServiceRequest struct {
	Name             string `json:"name" binding:"required,min=3"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	IconName         string `json:"icon_name"`
	IsFlagship       bool   `json:"is_flagship"`
	Status           string `json:"status" binding:"required,oneof=draft published"`
	// ThumbnailURL akan ditangkap oleh Service untuk disimpan ke media_assets
	ThumbnailURL     string `json:"thumbnail_url"` 
	ImpactPoints []string `json:"impact_points"`
	CTAText      string   `json:"cta_text"`
	CTALink      string   `json:"cta_link"`
}

// UpdateServiceRequest memvalidasi payload saat Admin memperbarui layanan
type UpdateServiceRequest struct {
	Name             string `json:"name" binding:"omitempty,min=3"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	IconName         string `json:"icon_name"`
	IsFlagship       *bool  `json:"is_flagship"` // Menggunakan pointer agar bisa mendeteksi nilai false
	Status           string `json:"status" binding:"omitempty,oneof=draft published"`
	ThumbnailURL     string `json:"thumbnail_url"`
	ImpactPoints []string `json:"impact_points"`
	CTAText      string   `json:"cta_text"`
	CTALink      string   `json:"cta_link"`
}

// ServiceQueryParams menangkap parameter URL untuk fitur Filter dan Pagination
type ServiceQueryParams struct {
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
	Search     string `form:"search"`
	Status     string `form:"status"`
	IsFlagship string `form:"is_flagship"` // Filter khusus untuk mencari layanan unggulan ("true" atau "false")
}