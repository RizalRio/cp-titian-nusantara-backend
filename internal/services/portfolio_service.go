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

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *PortfolioService) CreatePortfolio(req models.CreatePortfolioRequest, userID *uuid.UUID, ipAddress string) (*models.Portfolio, error) {
	var locations []models.PortfolioLocation
	for _, l := range req.Locations {
		locations = append(locations, models.PortfolioLocation{
			Name: l.Name,
			Lat:  l.Lat,
			Lng:  l.Lng,
		})
	}
	
	portfolio := models.Portfolio{
		Title:      req.Title,
		Slug:       GenerateSlug(req.Sector),
		Sector:     req.Sector,
		ShortStory: req.ShortStory,
		Impact:     req.Impact,
		Locations:  locations,
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

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction agar selaras
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Simpan ke Database
		if err := tx.Create(&portfolio).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (CREATE)
		LogActivity(tx, userID, "CREATE", "Portfolios", "Membuat Jejak Karya: "+portfolio.Title, ipAddress, nil, portfolio)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("judul jejak karya sudah digunakan")
		}
		return nil, err
	}

	// Kembalikan data lengkap beserta relasinya
	return s.repo.FindByID(portfolio.ID)
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *PortfolioService) UpdatePortfolio(id uuid.UUID, req models.UpdatePortfolioRequest, userID *uuid.UUID, ipAddress string) (*models.Portfolio, error) {
	portfolio, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("jejak karya tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Ambil snapshot data lama
	oldDataSnapshot := *portfolio

	// Update field dasar
	if req.Title != "" { portfolio.Title = req.Title }
	if req.Sector != "" { 
		portfolio.Sector = req.Sector 
		portfolio.Slug = GenerateSlug(req.Sector) 
	}
	if req.ShortStory != "" { portfolio.ShortStory = req.ShortStory }
	if req.Impact != "" { portfolio.Impact = req.Impact }
	if req.Locations != nil {
		var newLocs []models.PortfolioLocation
		for _, l := range req.Locations {
			newLocs = append(newLocs, models.PortfolioLocation{
				Name: l.Name,
				Lat:  l.Lat,
				Lng:  l.Lng,
			})
		}
		portfolio.Locations = newLocs // UPDATE LOKASI
	}
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
				if existing.FileURL != req.ThumbnailURL {
					tx.Delete(&existing) 
					tx.Create(&models.MediaAsset{ModelType: "Portfolio", ModelID: portfolio.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
				}
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Create(&models.MediaAsset{ModelType: "Portfolio", ModelID: portfolio.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
			}
		} else {
			var oldThumb models.MediaAsset
			if err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Portfolio", portfolio.ID, "thumbnail").First(&oldThumb).Error; err == nil {
				tx.Delete(&oldThumb) 
			}
		}

		// --- B. TANGANI GALERI (Delete and Insert) ---
		if req.GalleryURLs != nil {
			var oldGalleries []models.MediaAsset
			tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Portfolio", portfolio.ID, "gallery").Find(&oldGalleries)
			
			for _, m := range oldGalleries {
				tx.Delete(&m)
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

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Portfolios", "Memperbarui Jejak Karya: "+portfolio.Title, ipAddress, oldDataSnapshot, portfolio)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *PortfolioService) GetAllPortfolios(params models.PortfolioQueryParams) ([]models.Portfolio, int64, error) {
	return s.repo.FindAll(params)
}

func (s *PortfolioService) GetPortfolioByID(id uuid.UUID) (*models.Portfolio, error) {
	return s.repo.FindByID(id)
}

func (s *PortfolioService) GetPortfolioBySlug(slug string) (*models.Portfolio, error) {
	return s.repo.FindBySlug(slug)
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *PortfolioService) DeletePortfolio(id uuid.UUID, userID *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI LOG: Ambil data sebelum dihapus
	portfolioToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("jejak karya tidak ditemukan")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Cari dan hapus semua media terkait Program ini
		var media []models.MediaAsset
		// 🐛 FIX BUG: Sebelumnya tertulis "Project", diubah menjadi "Portfolio"
		tx.Where("model_type = ? AND model_id = ?", "Portfolio", id).Find(&media)
		for _, m := range media {
			tx.Delete(&m) // Memicu Hook Hapus Fisik
		}
		
		// Hapus record utama menggunakan tx
		if err := tx.Delete(&models.Portfolio{}, "id = ?", id).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "Portfolios", "Menghapus Jejak Karya: "+portfolioToDelete.Title, ipAddress, portfolioToDelete, nil)

		return nil
	})
}