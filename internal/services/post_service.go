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
		AuthorID:   authorID,
	}

	if req.Status == "published" {
		now := time.Now()
		post.PublishedAt = &now
	}

	var tags []models.Tag
	for _, tagID := range req.TagIDs {
		tags = append(tags, models.Tag{ID: tagID})
	}
	post.Tags = tags

	// 🌟 INJEKSI MEDIA ASSET (THUMBNAIL)
	// GORM otomatis akan mengisikan model_type="Post" dan model_id=(ID Artikel Baru)
	if req.ThumbnailURL != "" {
		post.Media = []models.MediaAsset{
			{
				MediaType: "thumbnail",
				FileURL:   req.ThumbnailURL,
			},
		}
	}

	if err := s.db.Create(&post).Error; err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("judul artikel sudah digunakan, silakan pilih yang lain")
		}
		return nil, err
	}

	return s.repo.FindByID(post.ID)
}

func (s *PostService) GetAllPosts(params models.PostQueryParams) ([]models.Post, int64, error) {
	return s.repo.FindAll(params)
}

func (s *PostService) GetPostByID(id uuid.UUID) (*models.Post, error) {
	return s.repo.FindByID(id)
}

// 🌟 GET POST BY SLUG
func (s *PostService) GetPostBySlug(slug string) (*models.Post, error) {
	// Panggil repository
	post, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan atau belum dipublikasikan")
	}
	return post, nil
}

func (s *PostService) UpdatePost(id uuid.UUID, req models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if req.Title != "" {
		post.Title = req.Title
		post.Slug = GenerateSlug(req.Title)
	}
	if req.CategoryID != uuid.Nil { post.CategoryID = req.CategoryID }
	if req.Content != "" { post.Content = req.Content }
	post.Excerpt = req.Excerpt
	
	if req.Status != "" && post.Status != req.Status {
		post.Status = req.Status
		if req.Status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	// 🔒 DATABASE TRANSACTION UNTUK UPDATE
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(post).Error; err != nil { return err }

		var newTags []models.Tag
		for _, tagID := range req.TagIDs {
			newTags = append(newTags, models.Tag{ID: tagID})
		}
		if err := tx.Model(post).Association("Tags").Replace(newTags); err != nil { return err }

		// 🌟 LOGIKA UPDATE THUMBNAIL POLIMORFIK
		if req.ThumbnailURL != "" {
			var existing models.MediaAsset
			err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Post", post.ID, "thumbnail").First(&existing).Error
			if err == nil {
				if existing.FileURL != req.ThumbnailURL {
					tx.Delete(&existing) // Memicu Hook Hapus Fisik
					tx.Create(&models.MediaAsset{ModelType: "Post", ModelID: post.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
				}
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Create(&models.MediaAsset{ModelType: "Post", ModelID: post.ID, MediaType: "thumbnail", FileURL: req.ThumbnailURL})
			}
		} else {
			var oldThumb models.MediaAsset
			if err := tx.Where("model_type = ? AND model_id = ? AND media_type = ?", "Post", post.ID, "thumbnail").First(&oldThumb).Error; err == nil {
				tx.Delete(&oldThumb) // Memicu Hook Hapus Fisik
			}
		}

		return nil
	})

	if err != nil { return nil, err }
	return s.repo.FindByID(id)
}

func (s *PostService) DeletePost(id uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Cari dan hapus media secara eksplisit agar Hook berjalan
		var media []models.MediaAsset
		tx.Where("model_type = ? AND model_id = ?", "Post", id).Find(&media)
		for _, m := range media {
			tx.Delete(&m) // Memicu Hook Hapus Fisik
		}
		
		// 2. Hapus relasi tags (jika tidak cascade)
		tx.Exec("DELETE FROM post_tags WHERE post_id = ?", id)
		
		// 3. Baru panggil repository untuk menghapus Post
		return s.repo.Delete(id)
	})
}