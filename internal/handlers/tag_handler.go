package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagHandler struct {
	service *services.TagService
}

func NewTagHandler(s *services.TagService) *TagHandler {
	return &TagHandler{service: s}
}

// 🌟 CREATE
func (h *TagHandler) Create(c *gin.Context) {
	var req models.CreateTagRequest
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

	tag, err := h.service.CreateTag(req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": tag})
}

// 🌟 READ ALL
func (h *TagHandler) GetAll(c *gin.Context) {
	tags, err := h.service.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": tags})
}

// 🌟 UPDATE
func (h *TagHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID tag tidak valid"})
		return
	}

	var req models.UpdateTagRequest
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

	tag, err := h.service.UpdateTag(id, req, userIDPtr, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": tag})
}

// 🌟 DELETE
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ID tag tidak valid"})
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

	if err := h.service.DeleteTag(id, userIDPtr, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal menghapus tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Tag berhasil dihapus"})
}