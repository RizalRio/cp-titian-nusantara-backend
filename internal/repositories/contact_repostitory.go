package repositories

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type ContactRepository struct {
	DB *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{DB: db}
}

func (r *ContactRepository) CreateMessage(msg *models.ContactMessage) error {
	return r.DB.Create(msg).Error
}

func (r *ContactRepository) CreateCollaboration(req *models.CollaborationRequest) error {
	return r.DB.Create(req).Error
}

// 🌟 ADMIN: Dapatkan Semua Pesan
func (r *ContactRepository) FindAllMessages(page, limit int, search string) ([]models.ContactMessage, int64, error) {
	var messages []models.ContactMessage
	var total int64
	query := r.DB.Model(&models.ContactMessage{})
	if search != "" {
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)
	err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&messages).Error
	return messages, total, err
}

// 🌟 ADMIN: Dapatkan Semua Kolaborasi
func (r *ContactRepository) FindAllCollaborations(page, limit int, search, status string) ([]models.CollaborationRequest, int64, error) {
	var collabs []models.CollaborationRequest
	var total int64
	query := r.DB.Model(&models.CollaborationRequest{})
	if search != "" {
		query = query.Where("LOWER(organization_name) LIKE ? OR LOWER(contact_person) LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&collabs).Error
	return collabs, total, err
}

// 🌟 ADMIN: Update Status & Hapus
func (r *ContactRepository) UpdateCollaborationStatus(id string, status string) error {
	return r.DB.Model(&models.CollaborationRequest{}).Where("id = ?", id).Update("status", status).Error
}
func (r *ContactRepository) MarkMessageAsRead(id string) error {
	return r.DB.Model(&models.ContactMessage{}).Where("id = ?", id).Update("is_read", true).Error
}
func (r *ContactRepository) DeleteCollaboration(id string) error {
	return r.DB.Delete(&models.CollaborationRequest{}, "id = ?", id).Error
}
func (r *ContactRepository) DeleteMessage(id string) error {
	return r.DB.Delete(&models.ContactMessage{}, "id = ?", id).Error
}