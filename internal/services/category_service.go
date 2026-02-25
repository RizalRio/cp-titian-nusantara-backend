package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	. "backend/pkg/utils"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(req models.CreateCategoryRequest) (*models.Category, error) {
	category := models.Category{
		Name: req.Name,
		Slug: GenerateSlug(req.Name),
	}

	if err := s.repo.Create(&category); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("kategori dengan nama ini sudah ada")
		}
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) UpdateCategory(id uuid.UUID, req models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil { return nil, errors.New("kategori tidak ditemukan") }

	category.Name = req.Name
	category.Slug = GenerateSlug(req.Name)

	if err := s.repo.Update(category); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("kategori dengan nama ini sudah ada")
		}
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	return s.repo.Delete(id)
}