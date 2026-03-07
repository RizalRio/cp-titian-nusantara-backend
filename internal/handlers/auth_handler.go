package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

// Login memproses request HTTP POST untuk login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// 1. Validasi input JSON (harus ada email & password)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format email tidak valid atau password kosong",
		})
		return
	}

	// 🌟 INJEKSI LOG: Ekstrak IP Address dari Context Gin
	ipAddress := c.ClientIP()

	// 2. Panggil Service Login beserta IP Address
	token, user, err := h.authService.Login(req, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// 3. Kirim Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login berhasil",
		"data": gin.H{
			"token": token,
			"user":  user, // PasswordHash otomatis tidak ikut terkirim karena `json:"-"` di model
		},
	})
}

// 🌟 FUNGSI BARU: LOGOUT
func (h *AuthHandler) Logout(c *gin.Context) {
	// 1. Ekstrak IP Address
	ipAddress := c.ClientIP()

	// 2. Ekstrak User ID dari Middleware JWT
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	// 3. Panggil Service Logout untuk merekam aktivitas
	if err := h.authService.Logout(userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mencatat log aktivitas keluar",
		})
		return
	}

	// 4. Kirim Response Sukses
	// Catatan: Penghapusan token asli (cookie/localStorage) tetap dilakukan oleh Frontend.
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil logout",
	})
}