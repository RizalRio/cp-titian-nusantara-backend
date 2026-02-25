package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	. "backend/pkg/utils"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type TagService struct {
	repo *repositories.TagRepository
}

func NewTagService(repo *repositories.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(req models.CreateTagRequest) (*models.Tag, error) {
	tag := models.Tag{
		Name: req.Name,
		Slug: GenerateSlug(req.Name),
	}

	if err := s.repo.Create(&tag); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("tag dengan nama ini sudah ada")
		}
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) GetAllTags() ([]models.Tag, error) {
	return s.repo.FindAll()
}

func (s *TagService) UpdateTag(id uuid.UUID, req models.UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil { return nil, errors.New("tag tidak ditemukan") }

	tag.Name = req.Name
	tag.Slug = GenerateSlug(req.Name)

	if err := s.repo.Update(tag); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("tag dengan nama ini sudah ada")
		}
		return nil, err
	}
	return tag, nil
}

func (s *TagService) DeleteTag(id uuid.UUID) error {
	return s.repo.Delete(id)
}