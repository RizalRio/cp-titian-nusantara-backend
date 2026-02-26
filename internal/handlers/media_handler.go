package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaHandler struct{}

func NewMediaHandler() *MediaHandler {
	return &MediaHandler{}
}

func (h *MediaHandler) UploadImage(c *gin.Context) {
	// 1. Tangkap file dari form-data dengan key "image"
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gagal mengunggah gambar. Pastikan key form-data adalah 'image'"})
		return
	}

	// 2. Validasi Ekstensi File
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format file tidak didukung. Gunakan PNG, JPG, WEBP, atau GIF"})
		return
	}

	// 3. Validasi Ukuran (Maks 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ukuran gambar maksimal 5MB"})
		return
	}

	// 4. Generate Nama File Unik (UUID) agar tidak ada file yang tertimpa
	newFileName := uuid.New().String() + ext
	uploadDir := "./uploads/images"

	// Buat folder jika belum ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal membuat direktori penyimpanan"})
		return
	}

	// 5. Simpan File ke Direktori Lokal
	savePath := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menyimpan gambar"})
		return
	}

	// 6. Kembalikan URL publik (Asumsi server berjalan di localhost:8080)
	// URL ini akan disimpan oleh Frontend ke database (sebagai thumbnail_url atau di dalam img src Quill)
	imageURL := fmt.Sprintf("http://localhost:8080/uploads/images/%s", newFileName)

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"url": imageURL,
		},
	})
}