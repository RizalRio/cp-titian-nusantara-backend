package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SettingService struct {
	repo *repositories.SettingRepository
}

func NewSettingService(repo *repositories.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

// 🌟 GET: Ubah banyak baris (DB) menjadi 1 DTO (Frontend)
func (s *SettingService) GetSettingsObject() (models.SiteSettingsDTO, error) {
	settings, err := s.repo.GetAllSettings()
	if err != nil {
		return models.SiteSettingsDTO{}, err
	}

	var dto models.SiteSettingsDTO
	for _, setting := range settings {
		switch setting.Key {
		case "site_name": dto.SiteName = setting.Value
		case "tagline": dto.Tagline = setting.Value
		case "description": dto.Description = setting.Value
		case "logo_url": dto.LogoURL = setting.Value
		case "favicon_url": dto.FaviconURL = setting.Value
		case "email": dto.Email = setting.Value
		case "phone": dto.Phone = setting.Value
		case "address": dto.Address = setting.Value
		case "instagram_url": dto.InstagramURL = setting.Value
		case "linkedin_url": dto.LinkedinURL = setting.Value
		case "youtube_url": dto.YoutubeURL = setting.Value
		}
	}
	return dto, nil
}

// 🌟 UPDATE: Ubah 1 DTO (Frontend) menjadi proses Upsert banyak baris (DB)
func (s *SettingService) UpdateSettings(req models.SiteSettingsDTO) error {
	// Peta konfigurasi: map[Key] struct{Value, Type}
	configMap := map[string]struct{ Val, Typ string }{
		"site_name":     {req.SiteName, "text"},
		"tagline":       {req.Tagline, "text"},
		"description":   {req.Description, "textarea"},
		"logo_url":      {req.LogoURL, "image"},
		"favicon_url":   {req.FaviconURL, "image"},
		"email":         {req.Email, "text"},
		"phone":         {req.Phone, "text"},
		"address":       {req.Address, "textarea"},
		"instagram_url": {req.InstagramURL, "url"},
		"linkedin_url":  {req.LinkedinURL, "url"},
		"youtube_url":   {req.YoutubeURL, "url"},
	}

	for key, data := range configMap {
		// Simpan atau perbarui tiap key ke database
		if err := s.repo.UpsertSetting(key, data.Val, data.Typ); err != nil {
			return err
		}
	}

	return nil
}