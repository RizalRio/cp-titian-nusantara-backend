package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	service *services.PostService
}

func NewPostHandler(s *services.PostService) *PostHandler {
	return &PostHandler{service: s}
}

// 🌟 CREATE: Dengan Keamanan AuthorID
func (h *PostHandler) Create(c *gin.Context) {
	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 🔒 SISTEM KEAMANAN: Ekstrak user_id dari JWT Middleware.
	// Kita berasumsi bahwa middleware login menyimpan ID user ke dalam context GIN.
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Sesi tidak valid atau tidak ada token JWT"})
		return
	}
	
	authorID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Format ID Pengguna tidak valid"})
		return
	}

	// 🌟 INJEKSI LOG: Ambil IP Address
	ipAddress := c.ClientIP()

	// Lemparkan data ke service, beserta ID pembuatnya dan data Log
	post, err := h.service.CreatePost(req, authorID, &authorID, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": post})
}

// 🌟 READ ALL: Menggunakan struktur respons dengan Meta Pagination
func (h *PostHandler) GetAll(c *gin.Context) {
	var params models.PostQueryParams
	
	// ShouldBindQuery akan menangkap URL seperti: /posts?page=1&limit=5&status=published
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Parameter filter tidak valid"})
		return
	}

	posts, total, err := h.service.GetAllPosts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data artikel"})
		return
	}

	// 💡 Menyusun JSON Response yang sangat ramah untuk Next.js (ada info total data untuk tombol "Next Page")
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   posts,
		"meta": gin.H{
			"total_data": total,
			"page":       params.Page,
			"limit":      params.Limit,
		},
	})
}

// 🌟 READ ONE: Mengambil detail satu artikel
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID artikel tidak valid"})
		return
	}

	post, err := h.service.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Artikel tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

// 🌟 GET BY SLUG (Untuk halaman detail pembaca)
func (h *PostHandler) GetBySlug(c *gin.Context) {
	// Menangkap parameter slug dari URL (contoh: /api/v1/posts/slug/membangun-desa)
	slug := c.Param("slug")
	
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Slug tidak boleh kosong"})
		return
	}

	post, err := h.service.GetPostBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success", 
		"data": post,
	})
}

// 🌟 UPDATE: Memperbarui artikel
func (h *PostHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID artikel tidak valid"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak UserID dan IP Address
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// 🌟 INJEKSI LOG: Panggil service dengan tambahan parameter
	post, err := h.service.UpdatePost(id, req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

// 🌟 DELETE: Menghapus artikel
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID artikel tidak valid"})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak UserID dan IP Address
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// 🌟 INJEKSI LOG: Panggil service dengan tambahan parameter
	if err := h.service.DeletePost(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus artikel"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Artikel berhasil dihapus"})
}