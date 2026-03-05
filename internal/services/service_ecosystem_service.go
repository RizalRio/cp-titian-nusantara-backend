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

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *ServiceEcosystemService) CreateService(req models.CreateServiceRequest, userID *uuid.UUID, ipAddress string) (*models.Service, error) {
	service := models.Service{
		Name:             req.Name,
		Slug:             GenerateSlug(req.Name), 
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		IconName:         req.IconName,
		IsFlagship:       req.IsFlagship,
		Status:           req.Status,
		ImpactPoints:     req.ImpactPoints, 
		CTAText:          req.CTAText,
		CTALink:          req.CTALink,
	}

	// Injeksi Media Asset (Thumbnail)
	if req.ThumbnailURL != "" {
		service.Media = []models.MediaAsset{
			{
				MediaType: "thumbnail",
				FileURL:   req.ThumbnailURL,
			},
		}
	}

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction agar selaras dengan LogActivity
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&service).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (CREATE)
		LogActivity(tx, userID, "CREATE", "Services", "Membuat Layanan: "+service.Name, ipAddress, nil, service)

		return nil
	})

	if err != nil {
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

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *ServiceEcosystemService) UpdateService(id uuid.UUID, req models.UpdateServiceRequest, userID *uuid.UUID, ipAddress string) (*models.Service, error) {
	service, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("layanan tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Ambil snapshot data lama untuk dicatat di log
	oldDataSnapshot := *service

	// Update Field Dasar
	if req.Name != "" {
		service.Name = req.Name
		service.Slug = GenerateSlug(req.Name)
	}
	if req.ShortDescription != "" { service.ShortDescription = req.ShortDescription }
	if req.Description != "" { service.Description = req.Description }
	if req.IconName != "" { service.IconName = req.IconName }
	if req.Status != "" { service.Status = req.Status }
	
	// Update Flagship
	if req.IsFlagship != nil {
		service.IsFlagship = *req.IsFlagship
	}

	if req.ImpactPoints != nil { service.ImpactPoints = req.ImpactPoints }
	if req.CTAText != "" { service.CTAText = req.CTAText }
	if req.CTALink != "" { service.CTALink = req.CTALink }

	// 🔒 DATABASE TRANSACTION
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan tabel services
		if err := tx.Save(service).Error; err != nil { return err }

		// 2. Logika Update Thumbnail Polimorfik
		if req.ThumbnailURL != "" {
			var existing models.MediaAsset
			err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Service", service.ID, "thumbnail").First(&existing).Error
			if err == nil {
				if existing.FileURL != req.ThumbnailURL {
					tx.Delete(&existing) // Memicu Hook
					tx.Create(&models.MediaAsset{ModelType: "Service", ModelID: service.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
				}
			} else {
				tx.Create(&models.MediaAsset{ModelType: "Service", ModelID: service.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
			}
		} else {
			var oldThumb models.MediaAsset
			if err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Service", service.ID, "thumbnail").First(&oldThumb).Error; err == nil {
				tx.Delete(&oldThumb) // Memicu Hook
			}
		}

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Services", "Memperbarui Layanan: "+service.Name, ipAddress, oldDataSnapshot, service)

		return nil
	})

	if err != nil { return nil, err }
	return s.repo.FindByID(id)
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *ServiceEcosystemService) DeleteService(id uuid.UUID, userID *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI LOG: Ambil data sebelum dihapus
	serviceToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("layanan tidak ditemukan")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Cari dan hapus semua media terkait Layanan ini
		var media []models.MediaAsset
		tx.Where("model_type = ? AND model_id = ?", "Service", id).Find(&media)
		for _, m := range media {
			tx.Delete(&m) // Memicu Hook Hapus Fisik
		}
		
		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "Services", "Menghapus Layanan: "+serviceToDelete.Name, ipAddress, serviceToDelete, nil)

		// Hapus record utama
		return s.repo.Delete(id)
	})
}