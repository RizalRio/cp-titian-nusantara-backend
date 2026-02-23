package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type SettingService struct {
	settingRepo *repositories.SettingRepository
}

func NewSettingService(repo *repositories.SettingRepository) *SettingService {
	return &SettingService{settingRepo: repo}
}

func (s *SettingService) GetAllSettings() ([]models.SiteSetting, error) {
	return s.settingRepo.GetAll()
}

func (s *SettingService) UpdateSettings(reqs []models.UpsertSettingRequest) error {
	var settings []models.SiteSetting

	// Mapping dari DTO ke Model
	for _, req := range reqs {
		settings = append(settings, models.SiteSetting{
			Key:   req.Key,
			Value: req.Value,
		})
	}

	return s.settingRepo.UpsertBatch(settings)
}