package models

// 🌟 DTO UNTUK FRONTEND (Satu Objek Datar)
type SiteSettingsDTO struct {
	SiteName     string `json:"site_name"`
	Tagline      string `json:"tagline"`
	Description  string `json:"description"`
	LogoURL      string `json:"logo_url"`
	FaviconURL   string `json:"favicon_url"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	InstagramURL string `json:"instagram_url"`
	LinkedinURL  string `json:"linkedin_url"`
	YoutubeURL   string `json:"youtube_url"`
}