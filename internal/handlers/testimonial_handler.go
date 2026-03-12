package handlers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestimonialHandler struct {
	service *services.TestimonialService
}

func NewTestimonialHandler(service *services.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{service: service}
}

// 🌟 GET ALL (Public) - Endpoint untuk ditarik oleh Frontend
func (h *TestimonialHandler) GetAll(c *gin.Context) {
	testimonialsData, err := h.service.GetAllTestimonials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", 
			"message": "Gagal mengambil data testimoni",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   testimonialsData,
	})
}