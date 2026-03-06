package services

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetDashboardStats() (models.DashboardStatsResponse, error) {
	var stats models.DashboardStatsResponse

	// Eksekusi COUNT() ke berbagai tabel secara langsung
	// Asumsi model artikel Anda bernama Post atau Article (sesuaikan jika berbeda)
	s.db.Model(&models.Post{}).Where("status = ?", "published").Count(&stats.TotalPosts) 
	
	s.db.Model(&models.Service{}).Count(&stats.TotalServices)
	
	s.db.Model(&models.Portfolio{}).Count(&stats.TotalPortfolios)
	
	// Hanya hitung pesan yang BELUM dibaca (is_read = false)
	s.db.Model(&models.ContactMessage{}).Where("is_read = ?", false).Count(&stats.UnreadMessages)

	return stats, nil
}