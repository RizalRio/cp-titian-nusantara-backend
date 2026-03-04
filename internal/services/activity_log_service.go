package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type ActivityLogService struct {
	repo *repositories.ActivityLogRepository
}

func NewActivityLogService(repo *repositories.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{repo: repo}
}

func (s *ActivityLogService) GetAllLogs(params models.ActivityLogQueryParams) ([]models.ActivityLog, int64, error) {
	if params.Page < 1 { params.Page = 1 }
	if params.Limit < 1 { params.Limit = 10 }
	return s.repo.FindAll(params)
}