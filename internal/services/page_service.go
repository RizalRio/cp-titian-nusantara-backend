package services

import (
	"errors"
	"time"

	"backend/internal/models"
	"backend/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PageService struct {
	pageRepo *repositories.PageRepository
	db       *gorm.DB // 🌟 INJEKSI: Diperlukan untuk Database Transaction pada Log
}

// 🌟 INJEKSI: Tambahkan db *gorm.DB
func NewPageService(repo *repositories.PageRepository, db *gorm.DB) *PageService {
	return &PageService{pageRepo: repo, db: db}
}

// 🌟 INJEKSI: Tambahkan parameter userIDPtr dan ipAddress
func (s *PageService) CreatePage(req models.CreatePageRequest, userID string, userIDPtr *uuid.UUID, ipAddress string) (*models.Page, error) {
	// Konversi input DTO ke Model Database
	page := &models.Page{
		Title:           req.Title,
		Slug:            req.Slug,
		TemplateName:    req.TemplateName,
		ContentJSON:     datatypes.JSON(req.ContentJSON),
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		Status:          req.Status,
		CreatedBy:       userID, // 🔒 Aman: Diambil dari token admin yang login
	}

	if req.Status == "published" {
		now := time.Now()
		page.PublishedAt = &now
	}

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Simpan ke database menggunakan tx
		if err := tx.Create(page).Error; err != nil {
			return errors.New("gagal menyimpan halaman, pastikan slug unik dan belum digunakan")
		}

		// 🌟 CATAT LOG AKTIVITAS (CREATE)
		LogActivity(tx, userIDPtr, "CREATE", "Pages", "Membuat Halaman: "+page.Title, ipAddress, nil, page)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (s *PageService) GetAllPages() ([]models.Page, error) {
	return s.pageRepo.FindAll()
}

func (s *PageService) GetPageByID(id string) (*models.Page, error) {
	return s.pageRepo.FindByID(id)
}

func (s *PageService) GetPageBySlug(slug string) (*models.Page, error) {
	return s.pageRepo.FindBySlug(slug)
}

// 🌟 INJEKSI: Tambahkan parameter userIDPtr dan ipAddress
func (s *PageService) UpdatePage(id string, req models.UpdatePageRequest, userID string, userIDPtr *uuid.UUID, ipAddress string) (*models.Page, error) {
	// 1. Cari data lamanya dulu
	page, err := s.pageRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("halaman tidak ditemukan")
	}

	// 🌟 INJEKSI: Ambil snapshot data lama untuk log
	oldDataSnapshot := *page

	// 2. Timpa data lama dengan data baru (jika diisi)
	if req.Title != "" { page.Title = req.Title }
	if req.Slug != "" { page.Slug = req.Slug }
	if req.TemplateName != "" { page.TemplateName = req.TemplateName }
	if req.MetaTitle != "" { page.MetaTitle = req.MetaTitle }
	if req.MetaDescription != "" { page.MetaDescription = req.MetaDescription }
	if req.ContentJSON != nil { page.ContentJSON = datatypes.JSON(req.ContentJSON) }

	// 3. Logika untuk status dan published_at
	if req.Status != "" {
		if page.Status != "published" && req.Status == "published" {
			now := time.Now()
			page.PublishedAt = &now
		}
		page.Status = req.Status
	}

	// 4. Catat siapa yang mengubah data ini
	page.UpdatedBy = &userID 

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Simpan pembaruan menggunakan tx
		if err := tx.Save(page).Error; err != nil {
			return errors.New("gagal memperbarui halaman, periksa kembali input Anda")
		}

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userIDPtr, "UPDATE", "Pages", "Memperbarui Halaman: "+page.Title, ipAddress, oldDataSnapshot, page)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return page, nil
}

// 🌟 INJEKSI: Tambahkan parameter userIDPtr dan ipAddress
func (s *PageService) DeletePage(id string, userIDPtr *uuid.UUID, ipAddress string) error {
	// 🌟 INJEKSI: Ambil data sebelum dihapus untuk dimasukkan ke log
	pageToDelete, err := s.pageRepo.FindByID(id)
	if err != nil {
		return errors.New("halaman tidak ditemukan")
	}

	// 🌟 INJEKSI LOG: Bungkus dengan tx.Transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Hapus menggunakan tx
		if err := tx.Delete(&models.Page{}, "id = ?", id).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userIDPtr, "DELETE", "Pages", "Menghapus Halaman: "+pageToDelete.Title, ipAddress, pageToDelete, nil)

		return nil
	})
}