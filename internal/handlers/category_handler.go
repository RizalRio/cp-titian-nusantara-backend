package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(s *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// 🌟 CREATE
func (h *CategoryHandler) Create(c *gin.Context) {
	var req models.CreateCategoryRequest

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

	// Kirim data yang sudah valid ke Service
	category, err := h.service.CreateCategory(req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": category})
}

// 🌟 READ ALL
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": categories})
}

// 🌟 UPDATE
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID kategori tidak valid"})
		return
	}

	var req models.UpdateCategoryRequest
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

	category, err := h.service.UpdateCategory(id, req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": category})
}

// 🌟 DELETE
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID kategori tidak valid"})
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

	if err := h.service.DeleteCategory(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus kategori (mungkin masih digunakan pada artikel)"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Kategori berhasil dihapus"})
}