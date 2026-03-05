package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingHandler struct {
	service *services.SettingService
}

func NewSettingHandler(service *services.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

// Publik & Admin GET API
func (h *SettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettingsObject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil pengaturan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": settings})
}

// Admin PUT API
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	var req models.SiteSettingsDTO
	
	// Bind JSON flat object dari Frontend ke Struct DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", 
			"message": "Format data tidak valid. Detail: " + err.Error(),
		})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak IP Address dan parsing UUID pengguna dari token JWT
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// 🌟 INJEKSI LOG: Lemparkan parameter ID dan IP ke service
	if err := h.service.UpdateSettings(req, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan pengaturan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pengaturan situs berhasil diperbarui"})
}