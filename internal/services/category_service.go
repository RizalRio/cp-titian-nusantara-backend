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

type CategoryService struct {
	repo *repositories.CategoryRepository
	db   *gorm.DB // 🌟 INJEKSI: Diperlukan untuk membungkus log dalam Transaction
}

// 🌟 INJEKSI: Tambahkan db *gorm.DB
func NewCategoryService(repo *repositories.CategoryRepository, db *gorm.DB) *CategoryService {
	return &CategoryService{repo: repo, db: db}
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *CategoryService) CreateCategory(req models.CreateCategoryRequest, userID *uuid.UUID, ipAddress string) (*models.Category, error) {
	category := models.Category{
		Name: req.Name,
		Slug: GenerateSlug(req.Name),
	}

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&category).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (CREATE)
		LogActivity(tx, userID, "CREATE", "Categories", "Membuat Kategori: "+category.Name, ipAddress, nil, category)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("kategori dengan nama ini sudah ada")
		}
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *CategoryService) UpdateCategory(id uuid.UUID, req models.UpdateCategoryRequest, userID *uuid.UUID, ipAddress string) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Ambil snapshot data lama
	oldDataSnapshot := *category

	category.Name = req.Name
	category.Slug = GenerateSlug(req.Name)

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(category).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Categories", "Memperbarui Kategori: "+category.Name, ipAddress, oldDataSnapshot, category)

		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("kategori dengan nama ini sudah ada")
		}
		return nil, err
	}

	return category, nil
}

// 🌟 INJEKSI LOG: Tambahkan parameter userID dan ipAddress
func (s *CategoryService) DeleteCategory(id uuid.UUID, userID *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI LOG: Ambil data kategori sebelum dihapus
	catToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("kategori tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Bungkus penghapusan dengan tx.Transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "Categories", "Menghapus Kategori: "+catToDelete.Name, ipAddress, catToDelete, nil)

		return nil
	})
}