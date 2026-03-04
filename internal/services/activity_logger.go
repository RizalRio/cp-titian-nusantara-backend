package services

import (
	"encoding/json"
	"log"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogActivity mencatat aktivitas ke tabel activity_logs.
// Wajib menggunakan tx (*gorm.DB) dari dalam blok tx.Transaction agar selaras dengan rollback/commit.
func LogActivity(tx *gorm.DB, userID *uuid.UUID, action, module, description, ipAddress string, oldData, newData interface{}) {
	var oldDataJSON, newDataJSON string

	if oldData != nil {
		if b, err := json.Marshal(oldData); err == nil {
			oldDataJSON = string(b)
		}
	}

	if newData != nil {
		if b, err := json.Marshal(newData); err == nil {
			newDataJSON = string(b)
		}
	}

	activityLog := models.ActivityLog{
		UserID:      userID,
		Action:      action,
		Module:      module,
		Description: description,
		OldData:     oldDataJSON,
		NewData:     newDataJSON,
		IPAddress:   ipAddress,
	}

	// Fail-Safe: Jika gagal log, jangan batalkan transaksi utama
	if err := tx.Create(&activityLog).Error; err != nil {
		log.Printf("⚠️ [Activity Log Error] Gagal mencatat aktivitas di modul %s: %v\n", module, err)
	}
}