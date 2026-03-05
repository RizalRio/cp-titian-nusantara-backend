package handlers

import (
	"math"
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PortfolioHandler struct {
	service *services.PortfolioService
}

func NewPortfolioHandler(service *services.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

// 🌟 CREATE (Admin Only)
func (h *PortfolioHandler) Create(c *gin.Context) {
	var req models.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// Lemparkan parameter log ke service
	portfolio, err := h.service.CreatePortfolio(req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": portfolio})
}

// 🌟 GET ALL (Public & Admin) - Mendukung Filter Sektor
func (h *PortfolioHandler) GetAll(c *gin.Context) {
	var params models.PortfolioQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	portfoliosData, totalData, err := h.service.GetAllPortfolios(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(params.Limit)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   portfoliosData,
		"meta": gin.H{
			"total_data":  totalData,
			"total_pages": totalPages,
			"page":        params.Page,
			"limit":       params.Limit,
		},
	})
}

// 🌟 GET BY ID (Admin / Umum)
func (h *PortfolioHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	portfolio, err := h.service.GetPortfolioByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Jejak Karya tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": portfolio})
}

// 🌟 GET BY SLUG (Khusus Halaman Publik SEO)
func (h *PortfolioHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Slug tidak boleh kosong"})
		return
	}

	portfolio, err := h.service.GetPortfolioBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Jejak Karya tidak ditemukan atau belum dipublikasikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": portfolio})
}

// 🌟 UPDATE (Admin Only)
func (h *PortfolioHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	var req models.UpdatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// Lemparkan parameter log ke service
	portfolio, err := h.service.UpdatePortfolio(id, req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": portfolio})
}

// 🌟 DELETE (Admin Only) - Soft Delete GORM
func (h *PortfolioHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// Lemparkan parameter log ke service
	if err := h.service.DeletePortfolio(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus jejak karya"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Jejak Karya berhasil dihapus"})
}