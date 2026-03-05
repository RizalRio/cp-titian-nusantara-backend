package services

import (
	"errors"

	"backend/internal/models"
	"backend/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactService struct {
	repo *repositories.ContactRepository
	db   *gorm.DB // 🌟 INJEKSI: Diperlukan untuk membungkus log dalam Transaction
}

// 🌟 INJEKSI: Tambahkan db *gorm.DB
func NewContactService(repo *repositories.ContactRepository, db *gorm.DB) *ContactService {
	return &ContactService{repo: repo, db: db}
}

func (s *ContactService) SubmitContactMessage(req models.CreateContactRequest) (*models.ContactMessage, error) {
	msg := &models.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	}
	err := s.repo.CreateMessage(msg)
	return msg, err
}

func (s *ContactService) SubmitCollaborationRequest(req models.CreateCollaborationRequest) (*models.CollaborationRequest, error) {
	collab := &models.CollaborationRequest{
		OrganizationName:  req.OrganizationName,
		ContactPerson:     req.ContactPerson,
		Email:             req.Email,
		Phone:             req.Phone,
		CollaborationType: req.CollaborationType,
		Message:           req.Message,
		ProposalFileURL:   req.ProposalFileURL,
		Status:            "pending",
	}
	err := s.repo.CreateCollaboration(collab)
	return collab, err
}

// 🌟 ADMIN: Dapatkan Semua Pesan Umum
func (s *ContactService) GetAllMessages(page, limit int, search string) ([]models.ContactMessage, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	return s.repo.FindAllMessages(page, limit, search)
}

// 🌟 ADMIN: Dapatkan Semua Kolaborasi
func (s *ContactService) GetAllCollaborations(page, limit int, search, status string) ([]models.CollaborationRequest, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	return s.repo.FindAllCollaborations(page, limit, search, status)
}

// 🌟 ADMIN: Ubah Status Kolaborasi (Diinjeksi Log)
func (s *ContactService) UpdateCollaborationStatus(id string, status string, userID *uuid.UUID, ipAddress string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var collab models.CollaborationRequest
		if err := tx.Where("id = ?", id).First(&collab).Error; err != nil {
			return errors.New("pengajuan kolaborasi tidak ditemukan")
		}

		// Ambil snapshot data lama
		oldDataSnapshot := collab

		// Update status
		if err := tx.Model(&collab).Update("status", status).Error; err != nil {
			return err
		}
		collab.Status = status // Update data baru untuk JSON NewData

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "Collaborations", "Mengubah status kolaborasi: "+collab.OrganizationName, ipAddress, oldDataSnapshot, collab)

		return nil
	})
}

// 🌟 ADMIN: Tandai Pesan Dibaca (Diinjeksi Log)
func (s *ContactService) MarkMessageAsRead(id string, userID *uuid.UUID, ipAddress string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var msg models.ContactMessage
		if err := tx.Where("id = ?", id).First(&msg).Error; err != nil {
			return errors.New("pesan tidak ditemukan")
		}

		// Ambil snapshot data lama
		oldDataSnapshot := msg

		// Update status baca
		if err := tx.Model(&msg).Update("is_read", true).Error; err != nil {
			return err
		}
		msg.IsRead = true // Update data baru untuk JSON NewData

		// 🌟 CATAT LOG AKTIVITAS (UPDATE)
		LogActivity(tx, userID, "UPDATE", "ContactMessages", "Menandai pesan telah dibaca dari: "+msg.Name, ipAddress, oldDataSnapshot, msg)

		return nil
	})
}

// 🌟 ADMIN: Hapus Pesan (Diinjeksi Log)
func (s *ContactService) DeleteMessage(id string, userID *uuid.UUID, ipAddress string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var msg models.ContactMessage
		if err := tx.Where("id = ?", id).First(&msg).Error; err != nil {
			return errors.New("pesan tidak ditemukan")
		}

		if err := tx.Delete(&msg).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "ContactMessages", "Menghapus pesan dari: "+msg.Name, ipAddress, msg, nil)

		return nil
	})
}

// 🌟 ADMIN: Hapus Kolaborasi (Diinjeksi Log)
func (s *ContactService) DeleteCollaboration(id string, userID *uuid.UUID, ipAddress string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var collab models.CollaborationRequest
		if err := tx.Where("id = ?", id).First(&collab).Error; err != nil {
			return errors.New("pengajuan kolaborasi tidak ditemukan")
		}

		if err := tx.Delete(&collab).Error; err != nil {
			return err
		}

		// 🌟 CATAT LOG AKTIVITAS (DELETE)
		LogActivity(tx, userID, "DELETE", "Collaborations", "Menghapus pengajuan kolaborasi dari: "+collab.OrganizationName, ipAddress, collab, nil)

		return nil
	})
}