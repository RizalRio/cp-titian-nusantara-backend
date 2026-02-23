package services

import (
	"errors"
	"time"

	"backend/internal/models"
	"backend/internal/repositories"

	"gorm.io/datatypes"
)

type PageService struct {
	pageRepo *repositories.PageRepository
}

func NewPageService(repo *repositories.PageRepository) *PageService {
	return &PageService{pageRepo: repo}
}

// CreatePage memproses pembuatan halaman baru
func (s *PageService) CreatePage(req models.CreatePageRequest, userID string) (*models.Page, error) {
	// Konversi input DTO ke Model Database
	page := &models.Page{
		Title:           req.Title,
		Slug:            req.Slug,
		TemplateName:    req.TemplateName,
		ContentJSON:     datatypes.JSON(req.ContentJSON), // Konversi aman dari json.RawMessage ke datatypes.JSON
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		Status:          req.Status,
		CreatedBy:       userID, // 🔒 Aman: Diambil dari token admin yang login
	}

	// Jika admin langsung memilih status "published", catat waktu rilisnya
	if req.Status == "published" {
		now := time.Now()
		page.PublishedAt = &now
	}

	// Simpan ke database
	err := s.pageRepo.Create(page)
	if err != nil {
		// Menangkap error jika Slug sudah dipakai (Unique Constraint)
		return nil, errors.New("gagal menyimpan halaman, pastikan slug unik dan belum digunakan")
	}

	return page, nil
}

// GetAllPages mengambil semua data halaman
func (s *PageService) GetAllPages() ([]models.Page, error) {
	return s.pageRepo.FindAll()
}

// GetPageByID mengambil detail satu halaman
func (s *PageService) GetPageByID(id string) (*models.Page, error) {
	return s.pageRepo.FindByID(id)
}

// UpdatePage memproses pembaruan data halaman
func (s *PageService) UpdatePage(id string, req models.UpdatePageRequest, userID string) (*models.Page, error) {
	// 1. Cari data lamanya dulu
	page, err := s.pageRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("halaman tidak ditemukan")
	}

	// 2. Timpa data lama dengan data baru (jika diisi)
	if req.Title != "" { page.Title = req.Title }
	if req.Slug != "" { page.Slug = req.Slug }
	if req.TemplateName != "" { page.TemplateName = req.TemplateName }
	if req.MetaTitle != "" { page.MetaTitle = req.MetaTitle }
	if req.MetaDescription != "" { page.MetaDescription = req.MetaDescription }
	if req.ContentJSON != nil { page.ContentJSON = datatypes.JSON(req.ContentJSON) }

	// 3. Logika untuk status dan published_at
	if req.Status != "" {
		// Jika sebelumnya draft dan sekarang di-publish
		if page.Status != "published" && req.Status == "published" {
			now := time.Now()
			page.PublishedAt = &now
		}
		page.Status = req.Status
	}

	// 4. Catat siapa yang mengubah data ini
	page.UpdatedBy = &userID // 🔒 Aman: Diambil dari token admin

	// 5. Simpan pembaruan
	err = s.pageRepo.Update(page)
	if err != nil {
		return nil, errors.New("gagal memperbarui halaman, periksa kembali input Anda")
	}

	return page, nil
}

// DeletePage menghapus halaman (Soft Delete)
func (s *PageService) DeletePage(id string) error {
	return s.pageRepo.Delete(id)
}