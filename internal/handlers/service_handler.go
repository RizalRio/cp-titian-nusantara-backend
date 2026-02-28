package handlers

import (
	"math"
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	service *services.ServiceEcosystemService
}

func NewServiceHandler(service *services.ServiceEcosystemService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

// 🌟 CREATE (Admin Only)
func (h *ServiceHandler) Create(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	service, err := h.service.CreateService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": service})
}

// 🌟 GET ALL (Public & Admin) dengan Pagination & Filter
func (h *ServiceHandler) GetAll(c *gin.Context) {
	var params models.ServiceQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	servicesData, totalData, err := h.service.GetAllServices(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(params.Limit)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   servicesData,
		"meta": gin.H{
			"total_data":  totalData,
			"total_pages": totalPages,
			"page":        params.Page,
			"limit":       params.Limit,
		},
	})
}

// 🌟 GET BY ID (Admin / Umum)
func (h *ServiceHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	service, err := h.service.GetServiceByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Layanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service})
}

// 🌟 GET BY SLUG (Khusus Halaman Publik SEO)
func (h *ServiceHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Slug tidak boleh kosong"})
		return
	}

	service, err := h.service.GetServiceBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Layanan tidak ditemukan atau belum dipublikasikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service})
}

// 🌟 UPDATE (Admin Only)
func (h *ServiceHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	var req models.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	service, err := h.service.UpdateService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service})
}

// 🌟 DELETE (Admin Only) - Soft Delete GORM
func (h *ServiceHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	if err := h.service.DeleteService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Layanan berhasil dihapus"})
}