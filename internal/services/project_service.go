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

type ProjectService struct {
	repo *repositories.ProjectRepository
	db   *gorm.DB
}

func NewProjectService(repo *repositories.ProjectRepository, db *gorm.DB) *ProjectService {
	return &ProjectService{repo: repo, db: db}
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *ProjectService) CreateProject(req models.CreateProjectRequest, userID *uuid.UUID, ipAddress string) (*models.Project, error) {
	project := models.Project{
		ServiceID:   req.ServiceID,
		Title:       req.Title,
		Slug:        GenerateSlug(req.Title),
		Summary:     req.Summary,
		Description: req.Description,
		Location:    req.Location,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      req.Status,
		IsFeatured:  req.IsFeatured,
	}

	var metrics []models.ProjectMetric
	for _, m := range req.Metrics {
		metrics = append(metrics, models.ProjectMetric{
			MetricKey:   m.MetricKey,
			MetricLabel: m.MetricLabel,
			MetricValue: m.MetricValue,
			MetricUnit:  m.MetricUnit,
			Order:       m.Order,
		})
	}
	if len(metrics) > 0 {
		project.Metrics = metrics
	}

	// 🌟 INJEKSI MEDIA (Thumbnail & Gallery)
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
		project.Media = mediaAssets
	}

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&project).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (CREATE)
		LogActivity(tx, userID, "CREATE", "Projects", "Membuat Project: "+project.Title, ipAddress, nil, project)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("judul project sudah digunakan")
		}
		return nil, err
	}

	return s.repo.FindByID(project.ID)
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *ProjectService) UpdateProject(id uuid.UUID, req models.UpdateProjectRequest, userID *uuid.UUID, ipAddress string) (*models.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil { return nil, errors.New("project tidak ditemukan") }

	// 🌟 INJEKSI LOG: Ambil snapshot data lama
	oldDataSnapshot := *project

	// Update field dasar
	if req.Title != "" { project.Title = req.Title; project.Slug = GenerateSlug(req.Title) }
	if req.ServiceID != uuid.Nil { project.ServiceID = req.ServiceID }
	if req.Summary != "" { project.Summary = req.Summary }
	if req.Description != "" { project.Description = req.Description }
	if req.Location != "" { project.Location = req.Location }
	if req.Status != "" { project.Status = req.Status }
	if req.IsFeatured != nil { project.IsFeatured = *req.IsFeatured }
	
	// Update tanggal (bisa di-set nil jika dihapus dari frontend)
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate

	// 🔒 DATABASE TRANSACTION UNTUK MEDIA POLIMORFIK
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(project).Error; err != nil { return err }

		if req.Metrics != nil {
			// Hapus semua metrik lama
			if err := tx.Where("project_id = ?", project.ID).Delete(&models.ProjectMetric{}).Error; err != nil {
				return err
			}
			
			// Masukkan metrik baru
			for _, m := range req.Metrics {
				newMetric := models.ProjectMetric{
					ProjectID:   project.ID,
					MetricKey:   m.MetricKey,
					MetricLabel: m.MetricLabel,
					MetricValue: m.MetricValue,
					MetricUnit:  m.MetricUnit,
					Order:       m.Order,
				}
				if err := tx.Create(&newMetric).Error; err != nil {
					return err
				}
			}
		}

		// 1. Tangani Thumbnail (Sama seperti Service/Post)
		if req.ThumbnailURL != "" {
			var existing models.MediaAsset
			err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Project", project.ID, "thumbnail").First(&existing).Error
			if err == nil {
				if existing.FileURL != req.ThumbnailURL {
					tx.Delete(&existing)
					tx.Create(&models.MediaAsset{ModelType: "Project", ModelID: project.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
				}
			} else {
				tx.Create(&models.MediaAsset{ModelType: "Project", ModelID: project.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
			}
		} else {
			var oldThumb models.MediaAsset
			if err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Project", project.ID, "thumbnail").First(&oldThumb).Error; err == nil {
				tx.Delete(&oldThumb)
			}
		}

		// 2. Tangani Gallery (Delete and Insert agar bersih)
		if req.GalleryURLs != nil {
			var oldGalleries []models.MediaAsset
			tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Project", project.ID, "gallery").Find(&oldGalleries)
			
			for _, m := range oldGalleries {
				tx.Delete(&m) // Memicu Hook
			}
			
			for _, url := range req.GalleryURLs {
				if url != "" {
					tx.Create(&models.MediaAsset{ModelType: "Project", ModelID: project.ID, MediaType: "gallery", FileURL: url})
				}
			}
		}

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Projects", "Memperbarui Project: "+project.Title, ipAddress, oldDataSnapshot, project)

		return nil
	})

	if err != nil { return nil, err }
	return s.repo.FindByID(id)
}

func (s *ProjectService) GetAllProjects(params models.ProjectQueryParams) ([]models.Project, int64, error) {
	return s.repo.FindAll(params)
}

func (s *ProjectService) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	return s.repo.FindByID(id)
}

func (s *ProjectService) GetProjectBySlug(slug string) (*models.Project, error) {   
	return s.repo.FindBySlug(slug)
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *ProjectService) DeleteProject(id uuid.UUID, userID *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI LOG: Ambil data sebelum dihapus
	projectToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Cari dan hapus semua media terkait Program ini
		var media []models.MediaAsset
		tx.Where("model_type = ? AND model_id = ?", "Project", id).Find(&media)
		for _, m := range media {
			tx.Delete(&m) // Memicu Hook Hapus Fisik
		}
		
		// Hapus record utama
		if err := tx.Delete(&models.Project{}, "id = ?", id).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "Projects", "Menghapus Project: "+projectToDelete.Title, ipAddress, projectToDelete, nil)

		return nil
	})
}