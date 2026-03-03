package repositories

import (
	"backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingRepository struct {
	DB *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{DB: db}
}

// Mengambil semua baris pengaturan
func (r *SettingRepository) GetAllSettings() ([]models.SiteSetting, error) {
	var settings []models.SiteSetting
	err := r.DB.Find(&settings).Error
	return settings, err
}

// Upsert: Update jika Key sudah ada, Insert jika belum ada
func (r *SettingRepository) UpsertSetting(key, value, settingType string) error {
	setting := models.SiteSetting{
		Key:   key,
		Value: value,
		Type:  settingType,
	}

	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}}, // Jika terjadi konflik pada 'key'
		DoUpdates: clause.AssignmentColumns([]string{"value", "type", "updated_at"}), // Update kolom ini
	}).Create(&setting).Error
}