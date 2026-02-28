package repositories

import (
	"strings"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

func (r *ProjectRepository) FindAll(params models.ProjectQueryParams) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	query := r.DB.Model(&models.Project{})

	// Filter
	if params.Status != "" { query = query.Where("status = ?", params.Status) }
	if params.IsFeatured != "" { query = query.Where("is_featured = ?", params.IsFeatured == "true") }
	if params.ServiceID != "" { query = query.Where("service_id = ?", params.ServiceID) }
	
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(location) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil { return nil, 0, err }

	offset := (params.Page - 1) * params.Limit
	err := query.
		Order("is_featured desc, created_at desc").
		Preload("Service"). // 🌟 Muat data layanan induk
		Preload("Media").   // 🌟 Muat thumbnail dan galeri
		Offset(offset).
		Limit(params.Limit).
		Find(&projects).Error

	return projects, total, err
}

func (r *ProjectRepository) FindByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.DB.Preload("Service").Preload("Media").First(&project, "id = ?", id).Error
	return &project, err
}

func (r *ProjectRepository) FindBySlug(slug string) (*models.Project, error) {
	var project models.Project
	err := r.DB.Where("status = ?", "published").Preload("Service").Preload("Media").First(&project, "slug = ?", slug).Error
	return &project, err
}

func (r *ProjectRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Project{}, "id = ?", id).Error
}