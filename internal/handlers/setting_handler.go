package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingService *services.SettingService
}

func NewSettingHandler(service *services.SettingService) *SettingHandler {
	return &SettingHandler{settingService: service}
}

// GetPublicSettings menangani GET /api/v1/settings (Tanpa Token)
func (h *SettingHandler) GetPublicSettings(c *gin.Context) {
	settings, err := h.settingService.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal memuat pengaturan situs"})
		return
	}

	// Format data menjadi objek (Map) agar Frontend Next.js lebih mudah membacanya
	// Contoh: { "contact_email": "halo@...", "footer_manifesto": "..." }
	formattedSettings := make(map[string]string)
	for _, s := range settings {
		formattedSettings[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": formattedSettings})
}

// UpdateSettings menangani PUT /api/v1/admin/settings (Wajib Token)
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	// Ingat, DTO kita adalah array of object JSON
	var req []models.UpsertSettingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format input tidak valid"})
		return
	}

	if err := h.settingService.UpdateSettings(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan pengaturan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pengaturan situs berhasil diperbarui"})
}