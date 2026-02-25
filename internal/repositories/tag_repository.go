package repositories

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct {
	DB *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{DB: db}
}

func (r *TagRepository) Create(tag *models.Tag) error {
	return r.DB.Create(tag).Error
}

func (r *TagRepository) FindAll() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.DB.Order("name asc").Find(&tags).Error
	return tags, err
}

func (r *TagRepository) FindByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.DB.First(&tag, "id = ?", id).Error
	return &tag, err
}

func (r *TagRepository) Update(tag *models.Tag) error {
	return r.DB.Save(tag).Error
}

func (r *TagRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Tag{}, "id = ?", id).Error
}