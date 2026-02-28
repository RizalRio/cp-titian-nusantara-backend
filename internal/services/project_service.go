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

func (s *ProjectService) CreateProject(req models.CreateProjectRequest) (*models.Project, error) {
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

	if err := s.db.Create(&project).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("judul project sudah digunakan")
		}
		return nil, err
	}

	return s.repo.FindByID(project.ID)
}

func (s *ProjectService) UpdateProject(id uuid.UUID, req models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil { return nil, errors.New("project tidak ditemukan") }

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

		// 1. Tangani Thumbnail (Sama seperti Service/Post)
		if req.ThumbnailURL != "" {
			var existing models.MediaAsset
			err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Project", project.ID, "thumbnail").First(&existing).Error
			if err == nil {
				existing.FileURL = req.ThumbnailURL
				if err := tx.Save(&existing).Error; err != nil { return err }
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Create(&models.MediaAsset{ModelType: "Project", ModelID: project.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
			}
		} else {
			tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Project", project.ID, "thumbnail").Delete(&models.MediaAsset{})
		}

		// 2. Tangani Gallery (Delete and Insert agar bersih)
		if req.GalleryURLs != nil {
			// Hapus semua galeri lama
			if err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Project", project.ID, "gallery").Delete(&models.MediaAsset{}).Error; err != nil {
				return err
			}
			
			// Masukkan galeri baru
			for _, url := range req.GalleryURLs {
				if url != "" {
					if err := tx.Create(&models.MediaAsset{ModelType: "Project", ModelID: project.ID, MediaType: "gallery", FileURL: url}).Error; err != nil {
						return err
					}
				}
			}
		}

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

func (s *ProjectService) DeleteProject(id uuid.UUID) error {
	return s.repo.Delete(id)
}
