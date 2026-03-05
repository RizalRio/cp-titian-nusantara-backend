package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContactHandler struct {
	service *services.ContactService
}

func NewContactHandler(service *services.ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

// (Publik) Submit Message & Collaboration tetap sama...
func (h *ContactHandler) SubmitMessage(c *gin.Context) {
	var req models.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	_, err := h.service.SubmitContactMessage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengirim pesan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Pesan berhasil dikirim"})
}

func (h *ContactHandler) SubmitCollaboration(c *gin.Context) {
	var req models.CreateCollaborationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	_, err := h.service.SubmitCollaborationRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengirim pengajuan kolaborasi"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Pengajuan kolaborasi berhasil dikirim"})
}

// 🌟 ADMIN: GET ALL CONTACT MESSAGES
func (h *ContactHandler) GetAllMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	messages, totalData, err := h.service.GetAllMessages(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data pesan"})
		return
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   messages,
		"meta": gin.H{
			"total_data":  totalData,
			"total_pages": totalPages,
			"page":        page,
			"limit":       limit,
		},
	})
}

// 🌟 ADMIN: GET ALL COLLABORATION REQUESTS
func (h *ContactHandler) GetAllCollaborations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")

	collabs, totalData, err := h.service.GetAllCollaborations(page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data kolaborasi"})
		return
	}

	totalPages := int(math.Ceil(float64(totalData) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   collabs,
		"meta": gin.H{
			"total_data":  totalData,
			"total_pages": totalPages,
			"page":        page,
			"limit":       limit,
		},
	})
}

// 🌟 ADMIN: MARK MESSAGE AS READ
func (h *ContactHandler) MarkMessageAsRead(c *gin.Context) {
	id := c.Param("id")

	// 🌟 INJEKSI LOG
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	if err := h.service.MarkMessageAsRead(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menandai pesan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pesan ditandai telah dibaca"})
}

// 🌟 ADMIN: UPDATE COLLABORATION STATUS
func (h *ContactHandler) UpdateCollaborationStatus(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Status string `json:"status" binding:"required,oneof=pending reviewed accepted rejected"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Status tidak valid"})
		return
	}

	// 🌟 INJEKSI LOG
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	if err := h.service.UpdateCollaborationStatus(id, req.Status, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengubah status kolaborasi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Status berhasil diubah"})
}

// 🌟 ADMIN: DELETE MESSAGE
func (h *ContactHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")

	// 🌟 INJEKSI LOG
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	if err := h.service.DeleteMessage(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus pesan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pesan berhasil dihapus"})
}

// 🌟 ADMIN: DELETE COLLABORATION
func (h *ContactHandler) DeleteCollaboration(c *gin.Context) {
	id := c.Param("id")

	// 🌟 INJEKSI LOG
	ipAddress := c.ClientIP()
	var userIDPtr *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			userIDPtr = &uid
		}
	}

	if err := h.service.DeleteCollaboration(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus pengajuan kolaborasi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Pengajuan kolaborasi berhasil dihapus"})
}