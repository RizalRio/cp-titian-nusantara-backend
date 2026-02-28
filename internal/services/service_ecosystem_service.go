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

type ServiceEcosystemService struct {
	repo *repositories.ServiceRepository
	db   *gorm.DB // Diperlukan untuk Database Transaction
}

func NewServiceEcosystemService(repo *repositories.ServiceRepository, db *gorm.DB) *ServiceEcosystemService {
	return &ServiceEcosystemService{repo: repo, db: db}
}

// 🌟 CREATE: Menyimpan Layanan dan Thumbnail Polimorfik
func (s *ServiceEcosystemService) CreateService(req models.CreateServiceRequest) (*models.Service, error) {
	service := models.Service{
		Name:             req.Name,
		Slug:             GenerateSlug(req.Name), // Asumsi fungsi GenerateSlug ada di utils.go
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		IconName:         req.IconName,
		IsFlagship:       req.IsFlagship,
		Status:           req.Status,
	}

	// Injeksi Media Asset (Thumbnail)
	// GORM otomatis akan mengisi ModelType="Service" dan ModelID=(ID Service Baru)
	if req.ThumbnailURL != "" {
		service.Media = []models.MediaAsset{
			{
				MediaType: "thumbnail",
				FileURL:   req.ThumbnailURL,
			},
		}
	}

	if err := s.db.Create(&service).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("nama layanan sudah digunakan, silakan pilih yang lain")
		}
		return nil, err
	}

	return s.repo.FindByID(service.ID)
}

func (s *ServiceEcosystemService) GetAllServices(params models.ServiceQueryParams) ([]models.Service, int64, error) {
	return s.repo.FindAll(params)
}

func (s *ServiceEcosystemService) GetServiceByID(id uuid.UUID) (*models.Service, error) {
	return s.repo.FindByID(id)
}

func (s *ServiceEcosystemService) GetServiceBySlug(slug string) (*models.Service, error) {
	return s.repo.FindBySlug(slug)
}

// 🌟 UPDATE: Menggunakan Transaction untuk Integritas Relasi Media
func (s *ServiceEcosystemService) UpdateService(id uuid.UUID, req models.UpdateServiceRequest) (*models.Service, error) {
	service, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("layanan tidak ditemukan")
	}

	// Update Field Dasar
	if req.Name != "" {
		service.Name = req.Name
		service.Slug = GenerateSlug(req.Name)
	}
	if req.ShortDescription != "" { service.ShortDescription = req.ShortDescription }
	if req.Description != "" { service.Description = req.Description }
	if req.IconName != "" { service.IconName = req.IconName }
	if req.Status != "" { service.Status = req.Status }
	
	// Update Flagship (Menggunakan pointer untuk mendeteksi false)
	if req.IsFlagship != nil {
		service.IsFlagship = *req.IsFlagship
	}

	// 🔒 DATABASE TRANSACTION
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan tabel services
		if err := tx.Save(service).Error; err != nil { return err }

		// 2. Logika Update Thumbnail Polimorfik
		if req.ThumbnailURL != "" {
			var existingThumbnail models.MediaAsset
			err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Service", service.ID, "thumbnail").First(&existingThumbnail).Error

			if err == nil {
				existingThumbnail.FileURL = req.ThumbnailURL
				if err := tx.Save(&existingThumbnail).Error; err != nil { return err }
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				newThumbnail := models.MediaAsset{
					ModelType: "Service",
					ModelID:   service.ID,
					MediaType: "thumbnail",
					FileURL:   req.ThumbnailURL,
				}
				if err := tx.Create(&newThumbnail).Error; err != nil { return err }
			} else {
				return err
			}
		} else {
			// Jika req.ThumbnailURL kosong, artinya Admin menghapus gambar. Kita hapus dari database.
			tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Service", service.ID, "thumbnail").Delete(&models.MediaAsset{})
		}

		return nil
	})

	if err != nil { return nil, err }
	return s.repo.FindByID(id)
}

func (s *ServiceEcosystemService) DeleteService(id uuid.UUID) error {
	return s.repo.Delete(id)
}