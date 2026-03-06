package handlers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetStats menangani permintaan GET /admin/dashboard-stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil statistik dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}