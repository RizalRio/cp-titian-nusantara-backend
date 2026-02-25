package repositories

import (
	"strings"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// FindAll memproses pencarian, filter, urutan, dan pagination
func (r *PostRepository) FindAll(params models.PostQueryParams) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	query := r.DB.Model(&models.Post{})

	// 1. Terapkan Filter Dinamis
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(excerpt) LIKE ?", searchTerm, searchTerm)
	}
	if params.TagSlug != "" {
		// Filter menggunakan JOIN ke tabel pivot post_tags
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.slug = ?", params.TagSlug)
	}

	// 2. Hitung Total Data (sebelum dilimit) untuk Pagination Frontend
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. Terapkan Sorting
	sortColumn := "created_at"
	if params.SortBy != "" { sortColumn = params.SortBy }
	sortOrder := "desc"
	if params.SortOrder == "asc" { sortOrder = "asc" }
	query = query.Order(sortColumn + " " + sortOrder)

	// 4. Terapkan Pagination
	offset := (params.Page - 1) * params.Limit

	// 5. Eksekusi dengan Preload (Eager Loading)
	err := query.
		Preload("Category").
		Preload("Author").
		Preload("Tags").
		Offset(offset).
		Limit(params.Limit).
		Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) FindByID(id uuid.UUID) (*models.Post, error) {
	var post models.Post
	err := r.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, "id = ?", id).Error
	return &post, err
}

func (r *PostRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Post{}, "id = ?", id).Error
}