package repositories

import (
	"strings"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioRepository struct {
	DB *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{DB: db}
}

func (r *PortfolioRepository) FindAll(params models.PortfolioQueryParams) ([]models.Portfolio, int64, error) {
	var portfolios []models.Portfolio
	var total int64

	query := r.DB.Model(&models.Portfolio{})

	if params.Status != "" { query = query.Where("status = ?", params.Status) }
	if params.Sector != "" { query = query.Where("sector = ?", params.Sector) }
	
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(location) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil { return nil, 0, err }

	offset := (params.Page - 1) * params.Limit
	err := query.Order("created_at desc").Preload("Media").Preload("Testimonials").Offset(offset).Limit(params.Limit).Find(&portfolios).Error

	return portfolios, total, err
}

func (r *PortfolioRepository) FindByID(id uuid.UUID) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	err := r.DB.Preload("Media").Preload("Testimonials").First(&portfolio, "id = ?", id).Error
	return &portfolio, err
}

func (r *PortfolioRepository) FindBySlug(slug string) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	err := r.DB.
		Preload("Media").
		Preload("Testimonials").
		Where("status = ?", "published").
		First(&portfolio, "slug = ?", slug).Error
	return &portfolio, err
}

func (r *PortfolioRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Portfolio{}, "id = ?", id).Error
}