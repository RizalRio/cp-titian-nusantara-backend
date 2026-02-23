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

// GetAll mengambil semua pengaturan (Bisa diakses publik oleh Frontend)
func (r *SettingRepository) GetAll() ([]models.SiteSetting, error) {
	var settings []models.SiteSetting
	err := r.DB.Find(&settings).Error
	return settings, err
}

// UpsertBatch memperbarui banyak pengaturan sekaligus
func (r *SettingRepository) UpsertBatch(settings []models.SiteSetting) error {
	// Jika "key" bentrok (sudah ada), maka update "value" dan "updated_at"
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&settings).Error
}