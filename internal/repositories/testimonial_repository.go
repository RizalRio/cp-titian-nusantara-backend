package repositories

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type TestimonialRepository struct {
	DB *gorm.DB
}

func NewTestimonialRepository(db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{DB: db}
}

// Menarik semua data testimoni, diurutkan dari yang terbaru
func (r *TestimonialRepository) FindAll() ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	
	// Mengambil semua data dari tabel testimonials, diurutkan berdasarkan waktu dibuat
	err := r.DB.Order("created_at desc").Find(&testimonials).Error
	
	return testimonials, err
}