package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	. "backend/pkg/utils"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagService struct {
	repo *repositories.TagRepository
	db   *gorm.DB // 🌟 INJEKSI: Diperlukan untuk membungkus log dalam Transaction
}

// 🌟 INJEKSI: Tambahkan db *gorm.DB
func NewTagService(repo *repositories.TagRepository, db *gorm.DB) *TagService {
	return &TagService{repo: repo, db: db}
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *TagService) CreateTag(req models.CreateTagRequest, userID *uuid.UUID, ipAddress string) (*models.Tag, error) {
	tag := models.Tag{
		Name: req.Name,
		Slug: GenerateSlug(req.Name),
	}

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Simpan tag menggunakan tx (bukan repo) agar berada dalam satu transaksi
		if err := tx.Create(&tag).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (CREATE)
		LogActivity(tx, userID, "CREATE", "Tags", "Membuat Tag: "+tag.Name, ipAddress, nil, tag)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("tag dengan nama ini sudah ada")
		}
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) GetAllTags() ([]models.Tag, error) {
	return s.repo.FindAll()
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *TagService) UpdateTag(id uuid.UUID, req models.UpdateTagRequest, userID *uuid.UUID, ipAddress string) (*models.Tag, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("tag tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Ambil snapshot data lama
	oldDataSnapshot := *tag

	tag.Name = req.Name
	tag.Slug = GenerateSlug(req.Name)

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(tag).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Tags", "Memperbarui Tag: "+tag.Name, ipAddress, oldDataSnapshot, tag)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("tag dengan nama ini sudah ada")
		}
		return nil, err
	}
	return tag, nil
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *TagService) DeleteTag(id uuid.UUID, userID *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI LOG: Ambil data tag sebelum dihapus
	tagToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("tag tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Bungkus penghapusan dengan tx.Transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Tag{}, "id = ?", id).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "Tags", "Menghapus Tag: "+tagToDelete.Name, ipAddress, tagToDelete, nil)

		return nil
	})
}