package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	pageService *services.PageService
}

func NewPageHandler(service *services.PageService) *PageHandler {
	return &PageHandler{pageService: service}
}

// Create menangani POST /admin/pages
func (h *PageHandler) Create(c *gin.Context) {
	var req models.CreatePageRequest

	// Validasi JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format input tidak valid. Pastikan semua field wajib diisi."})
		return
	}

	// 🔒 Mengambil User ID dari Middleware Auth
	userID := c.MustGet("user_id").(string)

	// Panggil Service
	page, err := h.pageService.CreatePage(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Halaman berhasil dibuat", "data": page})
}

// GetAll menangani GET /admin/pages
func (h *PageHandler) GetAll(c *gin.Context) {
	pages, err := h.pageService.GetAllPages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data halaman"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": pages})
}

// GetByID menangani GET /admin/pages/:id
func (h *PageHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	page, err := h.pageService.GetPageByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Halaman tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": page})
}

// Update menangani PUT /admin/pages/:id
func (h *PageHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdatePageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format input tidak valid"})
		return
	}

	// 🔒 Mengambil User ID dari Middleware Auth
	userID := c.MustGet("user_id").(string)

	page, err := h.pageService.UpdatePage(id, req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Halaman berhasil diperbarui", "data": page})
}

// Delete menangani DELETE /admin/pages/:id
func (h *PageHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.pageService.DeletePage(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus halaman"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Halaman berhasil dihapus"})
}