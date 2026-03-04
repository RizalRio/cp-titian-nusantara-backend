package repositories

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type ActivityLogRepository struct {
	DB *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{DB: db}
}

// Fungsi internal untuk mencatat log dari service lain
func (r *ActivityLogRepository) CreateLog(log *models.ActivityLog) error {
	return r.DB.Create(log).Error
}

// Dapatkan semua log untuk Admin
func (r *ActivityLogRepository) FindAll(params models.ActivityLogQueryParams) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	query := r.DB.Model(&models.ActivityLog{}).Preload("User")

	if params.Search != "" {
		query = query.Where("LOWER(description) LIKE ? OR LOWER(ip_address) LIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}
	if params.Module != "" && params.Module != "all" {
		query = query.Where("module = ?", params.Module)
	}
	if params.Action != "" && params.Action != "all" {
		query = query.Where("action = ?", params.Action)
	}

	query.Count(&total)

	err := query.Order("created_at desc").
		Offset((params.Page - 1) * params.Limit).
		Limit(params.Limit).
		Find(&logs).Error

	return logs, total, err
}