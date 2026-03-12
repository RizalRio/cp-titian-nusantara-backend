package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
)

type TestimonialService struct {
	repo *repositories.TestimonialRepository
}

func NewTestimonialService(repo *repositories.TestimonialRepository) *TestimonialService {
	return &TestimonialService{repo: repo}
}

func (s *TestimonialService) GetAllTestimonials() ([]models.Testimonial, error) {
	return s.repo.FindAll()
}