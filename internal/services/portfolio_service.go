package services

import (
	"errors"
	"strings"

	"backend/internal/models"
	"backend/internal/repositories"
	. "backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortfolioService struct {
	repo *repositories.PortfolioRepository
	db   *gorm.DB
}

func NewPortfolioService(repo *repositories.PortfolioRepository, db *gorm.DB) *PortfolioService {
	return &PortfolioService{repo: repo, db: db}
}

// 🌟 CREATE PORTFOLIO (Jejak Karya)
func (s *PortfolioService) CreatePortfolio(req models.CreatePortfolioRequest) (*models.Portfolio, error) {
	portfolio := models.Portfolio{
		Title:      req.Title,
		Slug:       GenerateSlug(req.Sector),
		Sector:     req.Sector,
		ShortStory: req.ShortStory,
		Impact:     req.Impact,
		Location:   req.Location,
		Status:     req.Status,
	}

	// 1. Tangani Media (Thumbnail & Gallery)
	var mediaAssets []models.MediaAsset
	if req.ThumbnailURL != "" {
		mediaAssets = append(mediaAssets, models.MediaAsset{
			MediaType: "thumbnail",
			FileURL:   req.ThumbnailURL,
		})
	}
	for _, url := range req.GalleryURLs {
		if url != "" {
			mediaAssets = append(mediaAssets, models.MediaAsset{
				MediaType: "gallery",
				FileURL:   url,
			})
		}
	}
	if len(mediaAssets) > 0 {
		portfolio.Media = mediaAssets
	}

	// 2. Tangani Testimonial (One-to-Many)
	var testimonials []models.Testimonial
	for _, t := range req.Testimonials {
		testimonials = append(testimonials, models.Testimonial{
			AuthorName: t.AuthorName,
			AuthorRole: t.AuthorRole,
			Content:    t.Content,
			AvatarURL:  t.AvatarURL,
			Order:      t.Order,
		})
	}
	if len(testimonials) > 0 {
		portfolio.Testimonials = testimonials
	}

	// Simpan ke Database
	if err := s.db.Create(&portfolio).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("judul jejak karya sudah digunakan")
		}
		return nil, err
	}

	// Kembalikan data lengkap beserta relasinya
	return s.repo.FindByID(portfolio.ID)
}

// 🌟 UPDATE PORTFOLIO (Jejak Karya)
func (s *PortfolioService) UpdatePortfolio(id uuid.UUID, req models.UpdatePortfolioRequest) (*models.Portfolio, error) {
	portfolio, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("jejak karya tidak ditemukan")
	}

	// Update field dasar
	if req.Title != "" { portfolio.Title = req.Title }
	if req.Sector != "" { 
		portfolio.Sector = req.Sector 
		portfolio.Slug = GenerateSlug(req.Sector) 
	}
	if req.ShortStory != "" { portfolio.ShortStory = req.ShortStory }
	if req.Impact != "" { portfolio.Impact = req.Impact }
	if req.Location != "" { portfolio.Location = req.Location }
	if req.Status != "" { portfolio.Status = req.Status }

	// 🔒 DATABASE TRANSACTION UNTUK MEDIA & TESTIMONI
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Simpan perubahan dasar Portfolio
		if err := tx.Save(portfolio).Error; err != nil {
			return err
		}

		// --- A. TANGANI MEDIA THUMBNAIL ---
		if req.ThumbnailURL != "" {
			var existing models.MediaAsset
			err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Portfolio", portfolio.ID, "thumbnail").First(&existing).Error
			if err == nil {
				existing.FileURL = req.ThumbnailURL
				if err := tx.Save(&existing).Error; err != nil { return err }
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Create(&models.MediaAsset{ModelType: "Portfolio", ModelID: portfolio.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
			}
		} else {
			tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Portfolio", portfolio.ID, "thumbnail").Delete(&models.MediaAsset{})
		}

		// --- B. TANGANI GALERI (Delete and Insert) ---
		if req.GalleryURLs != nil {
			if err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Portfolio", portfolio.ID, "gallery").Delete(&models.MediaAsset{}).Error; err != nil {
				return err
			}
			for _, url := range req.GalleryURLs {
				if url != "" {
					if err := tx.Create(&models.MediaAsset{ModelType: "Portfolio", ModelID: portfolio.ID, MediaType: "gallery", FileURL: url}).Error; err != nil {
						return err
					}
				}
			}
		}

		// --- C. TANGANI TESTIMONI (Delete and Insert) ---
		if req.Testimonials != nil {
			// Hapus semua testimoni lama
			if err := tx.Where("portfolio_id = ?", portfolio.ID).Delete(&models.Testimonial{}).Error; err != nil {
				return err
			}
			// Masukkan testimoni baru
			for _, t := range req.Testimonials {
				newTestimonial := models.Testimonial{
					PortfolioID: portfolio.ID,
					AuthorName:  t.AuthorName,
					AuthorRole:  t.AuthorRole,
					Content:     t.Content,
					AvatarURL:   t.AvatarURL,
					Order:       t.Order,
				}
				if err := tx.Create(&newTestimonial).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

// 🌟 GET ALL PORTFOLIOS
func (s *PortfolioService) GetAllPortfolios(params models.PortfolioQueryParams) ([]models.Portfolio, int64, error) {
	return s.repo.FindAll(params)
}

// 🌟 GET PORTFOLIO BY ID
func (s *PortfolioService) GetPortfolioByID(id uuid.UUID) (*models.Portfolio, error) {
	return s.repo.FindByID(id)
}

// 🌟 GET PORTFOLIO BY SLUG (Untuk Publik)
func (s *PortfolioService) GetPortfolioBySlug(slug string) (*models.Portfolio, error) {
	return s.repo.FindBySlug(slug)
}

// 🌟 DELETE PORTFOLIO
func (s *PortfolioService) DeletePortfolio(id uuid.UUID) error {
	return s.repo.Delete(id)
}