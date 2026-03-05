package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format input tidak valid. Pastikan semua field wajib diisi."})
		return
	}

	userID := c.MustGet("user_id").(string)

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if uid, err := uuid.Parse(userID); err == nil {
		userIDPtr = &uid
	}

	// Panggil Service dengan parameter tambahan
	page, err := h.pageService.CreatePage(req, userID, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Halaman berhasil dibuat", "data": page})
}

func (h *PageHandler) GetAll(c *gin.Context) {
	pages, err := h.pageService.GetAllPages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data halaman"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": pages})
}

func (h *PageHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	page, err := h.pageService.GetPageByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Halaman tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": page})
}

func (h *PageHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	page, err := h.pageService.GetPageBySlug(slug)
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

	userID := c.MustGet("user_id").(string)

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if uid, err := uuid.Parse(userID); err == nil {
		userIDPtr = &uid
	}

	page, err := h.pageService.UpdatePage(id, req, userID, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Halaman berhasil diperbarui", "data": page})
}

// Delete menangani DELETE /admin/pages/:id
func (h *PageHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	if err := h.pageService.DeletePage(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus halaman"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Halaman berhasil dihapus"})
}