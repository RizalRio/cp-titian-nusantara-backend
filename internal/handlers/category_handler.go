package handlers

import (
	"net/http"

	"backend/internal/models" // Sesuaikan dengan nama module golang kamu
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	service *services.CategoryService
}

// Constructor untuk menginisialisasi handler dengan service-nya
func NewCategoryHandler(s *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// 🌟 CREATE: Menerima JSON dari frontend Next.js untuk membuat kategori baru
func (h *CategoryHandler) Create(c *gin.Context) {
	var req models.CreateCategoryRequest
	
	// ShouldBindJSON akan mencocokkan JSON body dengan DTO. 
	// Jika ada field yang kurang (misal nama terlalu pendek), GIN otomatis menolaknya.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Kirim data yang sudah valid ke Service
	category, err := h.service.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Kembalikan status 201 (Created) beserta data kategorinya
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": category})
}

// 🌟 READ ALL: Mengirimkan daftar kategori ke frontend
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data kategori"})
		return
	}
	// Mengembalikan status 200 (OK)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": categories})
}

// 🌟 UPDATE: Memperbarui nama kategori
func (h *CategoryHandler) Update(c *gin.Context) {
	// Tangkap "id" dari URL (contoh: /api/categories/123e4567-...)
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

	category, err := h.service.UpdateCategory(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": category})
}

// 🌟 DELETE: Menghapus kategori
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID kategori tidak valid"})
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus kategori (mungkin masih digunakan pada artikel)"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Kategori berhasil dihapus"})
}