package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type ContactService struct {
	repo *repositories.ContactRepository
}

func NewContactService(repo *repositories.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
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

// 🌟 ADMIN: Ubah Status Kolaborasi
func (s *ContactService) UpdateCollaborationStatus(id string, status string) error {
	return s.repo.UpdateCollaborationStatus(id, status)
}

// 🌟 ADMIN: Tandai Pesan Dibaca
func (s *ContactService) MarkMessageAsRead(id string) error {
	return s.repo.MarkMessageAsRead(id)
}

// 🌟 ADMIN: Hapus Pesan
func (s *ContactService) DeleteMessage(id string) error {
	return s.repo.DeleteMessage(id)
}

// 🌟 ADMIN: Hapus Kolaborasi
func (s *ContactService) DeleteCollaboration(id string) error {
	return s.repo.DeleteCollaboration(id)
}