package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// 2. Panggil Service Login
	token, user, err := h.authService.Login(req)
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