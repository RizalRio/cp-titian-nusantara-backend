package repositories

import (
	"strings"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	DB *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{DB: db}
}

// 🌟 READ ALL: Dengan Filter, Search, Pagination, dan Preload Media
func (r *ServiceRepository) FindAll(params models.ServiceQueryParams) ([]models.Service, int64, error) {
	var services []models.Service
	var total int64

	query := r.DB.Model(&models.Service{})

	// 1. Filter Dinamis
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.IsFlagship != "" {
		// Konversi string "true"/"false" dari URL query ke boolean SQL
		isFlagship := params.IsFlagship == "true"
		query = query.Where("is_flagship = ?", isFlagship)
	}
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(short_description) LIKE ?", searchTerm, searchTerm)
	}

	// 2. Hitung Total Data untuk Pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. Terapkan Pagination & Urutan (Layanan unggulan tampil duluan, lalu yang terbaru)
	offset := (params.Page - 1) * params.Limit
	err := query.
		Order("is_flagship desc, created_at desc").
		Preload("Media"). // Eager loading thumbnail
		Offset(offset).
		Limit(params.Limit).
		Find(&services).Error

	return services, total, err
}

func (r *ServiceRepository) FindByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	err := r.DB.Preload("Media").First(&service, "id = ?", id).Error
	return &service, err
}

func (r *ServiceRepository) FindBySlug(slug string) (*models.Service, error) {
	var service models.Service
	err := r.DB.Where("status = ?", "published").Preload("Media").First(&service, "slug = ?", slug).Error
	return &service, err
}

func (r *ServiceRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Service{}, "id = ?", id).Error
}