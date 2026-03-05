package services

import (
	"backend/internal/models"
	"backend/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SettingService struct {
	repo *repositories.SettingRepository
	db   *gorm.DB // 🌟 INJEKSI: Diperlukan untuk membungkus log dalam Transaction
}

// 🌟 INJEKSI: Tambahkan db *gorm.DB
func NewSettingService(repo *repositories.SettingRepository, db *gorm.DB) *SettingService {
	return &SettingService{repo: repo, db: db}
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

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *SettingService) UpdateSettings(req models.SiteSettingsDTO, userID *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI LOG: Ambil snapshot data lama sebelum diperbarui
	oldDataSnapshot, _ := s.GetSettingsObject()

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

	// 🌟 INJEKSI LOG: Bungkus proses Upsert di dalam tx.Transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		for key, data := range configMap {
			// Simpan atau perbarui tiap key ke database
			if err := s.repo.UpsertSetting(key, data.Val, data.Typ); err != nil {
				return err
			}
		}

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Settings", "Memperbarui konfigurasi identitas dan pengaturan situs", ipAddress, oldDataSnapshot, req)

		return nil
	})
}