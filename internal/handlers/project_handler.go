package handlers

import (
	"math"
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// 🌟 CREATE (Admin Only)
func (h *ProjectHandler) Create(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	project, err := h.service.CreateProject(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": project})
}

// 🌟 GET ALL (Public & Admin) dengan Filter Ekstensif
func (h *ProjectHandler) GetAll(c *gin.Context) {
	var params models.ProjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	projectsData, totalData, err := h.service.GetAllProjects(params) // Asumsi fungsi ini sudah Anda tambahkan di Service
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(params.Limit)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   projectsData,
		"meta": gin.H{
			"total_data":  totalData,
			"total_pages": totalPages,
			"page":        params.Page,
			"limit":       params.Limit,
		},
	})
}

// 🌟 GET BY ID (Admin / Umum)
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	project, err := h.service.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Project tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": project})
}

// 🌟 GET BY SLUG (Halaman Publik SEO)
func (h *ProjectHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Slug tidak boleh kosong"})
		return
	}

	project, err := h.service.GetProjectBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Project tidak ditemukan atau belum dipublikasikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": project})
}

// 🌟 UPDATE (Admin Only)
func (h *ProjectHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	project, err := h.service.UpdateProject(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": project})
}

// 🌟 DELETE (Admin Only) - Soft Delete
func (h *ProjectHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format ID tidak valid"})
		return
	}

	if err := h.service.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Project berhasil dihapus"})
}