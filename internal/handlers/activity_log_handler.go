package handlers

import (
	"math"
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ActivityLogHandler struct {
	service *services.ActivityLogService
}

func NewActivityLogHandler(service *services.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

func (h *ActivityLogHandler) GetAllLogs(c *gin.Context) {
	var params models.ActivityLogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Parameter tidak valid"})
		return
	}

	logs, totalData, err := h.service.GetAllLogs(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil log aktivitas"})
		return
	}

	// Cegah limit 0 untuk pembagian
	if params.Limit == 0 { params.Limit = 10 }
	totalPages := int(math.Ceil(float64(totalData) / float64(params.Limit)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   logs,
		"meta": gin.H{
			"total_data":  totalData,
			"total_pages": totalPages,
			"page":        params.Page,
			"limit":       params.Limit,
		},
	})
}