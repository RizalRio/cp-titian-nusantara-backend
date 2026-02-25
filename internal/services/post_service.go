package services

import (
	"errors"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repositories"
	. "backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	repo *repositories.PostRepository
	db   *gorm.DB // Diperlukan untuk Database Transaction pada relasi Tags
}

func NewPostService(repo *repositories.PostRepository, db *gorm.DB) *PostService {
	return &PostService{repo: repo, db: db}
}

func (s *PostService) CreatePost(req models.CreatePostRequest, authorID uuid.UUID) (*models.Post, error) {
	post := models.Post{
		Title:      req.Title,
		Slug:       GenerateSlug(req.Title),
		CategoryID: req.CategoryID,
		Excerpt:    req.Excerpt,
		Content:    req.Content,
		Status:     req.Status,
		AuthorID:   authorID, // 🔒 Aman dari spoofing frontend
	}

	if req.Status == "published" {
		now := time.Now()
		post.PublishedAt = &now
	}

	// 🌟 Membangun relasi Tag secara dinamis
	var tags []models.Tag
	for _, tagID := range req.TagIDs {
		tags = append(tags, models.Tag{ID: tagID})
	}
	post.Tags = tags

	// Simpan ke DB (GORM otomatis akan insert ke tabel posts dan pivot post_tags)
	if err := s.db.Create(&post).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("judul artikel sudah digunakan, silakan pilih yang lain")
		}
		return nil, err
	}

	// Mengembalikan data post lengkap dengan preload Category & Tags
	return s.repo.FindByID(post.ID)
}

func (s *PostService) GetAllPosts(params models.PostQueryParams) ([]models.Post, int64, error) {
	return s.repo.FindAll(params)
}

func (s *PostService) GetPostByID(id uuid.UUID) (*models.Post, error) {
	return s.repo.FindByID(id)
}

func (s *PostService) UpdatePost(id uuid.UUID, req models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	// Update Field Dasar
	if req.Title != "" {
		post.Title = req.Title
		post.Slug = GenerateSlug(req.Title)
	}
	if req.CategoryID != uuid.Nil { post.CategoryID = req.CategoryID }
	if req.Content != "" { post.Content = req.Content }
	post.Excerpt = req.Excerpt
	
	// Cek Status Publikasi
	if req.Status != "" && post.Status != req.Status {
		post.Status = req.Status
		if req.Status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	// 🔒 DATABASE TRANSACTION: Untuk update artikel & relasi tag dengan aman
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan perubahan pada tabel posts
		if err := tx.Save(post).Error; err != nil {
			return err
		}

		// 2. Siapkan array Tag baru
		var newTags []models.Tag
		for _, tagID := range req.TagIDs {
			newTags = append(newTags, models.Tag{ID: tagID})
		}
		
		// 3. .Replace() akan otomatis menghapus tag lama di pivot dan memasukkan yang baru
		if err := tx.Model(post).Association("Tags").Replace(newTags); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Kembalikan data terbaru setelah diupdate
	return s.repo.FindByID(id)
}

func (s *PostService) DeletePost(id uuid.UUID) error {
	return s.repo.Delete(id)
}