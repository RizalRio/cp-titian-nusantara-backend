package repositories

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type PageRepository struct {
	DB *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{DB: db}
}

// Create menyimpan halaman baru ke database
func (r *PageRepository) Create(page *models.Page) error {
	return r.DB.Create(page).Error
}

// FindAll mengambil semua halaman (untuk tabel di dashboard admin)
func (r *PageRepository) FindAll() ([]models.Page, error) {
	var pages []models.Page
	err := r.DB.Order("created_at desc").Find(&pages).Error
	return pages, err
}

// FindByID mencari satu halaman untuk diedit
func (r *PageRepository) FindByID(id string) (*models.Page, error) {
	var page models.Page
	err := r.DB.Where("id = ?", id).First(&page).Error
	return &page, err
}

// Update menyimpan perubahan data halaman
func (r *PageRepository) Update(page *models.Page) error {
	return r.DB.Save(page).Error
}

// Delete melakukan soft-delete pada halaman
func (r *PageRepository) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&models.Page{}).Error
}