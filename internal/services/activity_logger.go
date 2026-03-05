package services

import (
	"encoding/json"
	"log"

	"backend/internal/models" // Sesuaikan path jika berbeda

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogActivity mencatat aktivitas ke tabel activity_logs.
func LogActivity(tx *gorm.DB, userID *uuid.UUID, action, module, description, ipAddress string, oldData, newData interface{}) {
	
	// 🌟 SOLUSI ANTI-PELURU: Gunakan objek JSON kosong "{}" alih-alih "null" atau ""
	oldDataJSON := "{}"
	newDataJSON := "{}"

	// Parsing Data Lama
	if oldData != nil {
		if b, err := json.Marshal(oldData); err == nil {
			str := string(b)
			// Pastikan hasil marshal bukan string kosong atau sekadar "null"
			if str != "" && str != "null" {
				oldDataJSON = str
			}
		}
	}

	// Parsing Data Baru
	if newData != nil {
		if b, err := json.Marshal(newData); err == nil {
			str := string(b)
			if str != "" && str != "null" {
				newDataJSON = str
			}
		}
	}

	activityLog := models.ActivityLog{
		UserID:      userID,
		Action:      action,
		Module:      module,
		Description: description,
		OldData:     oldDataJSON, // Pasti berisi minimal "{}"
		NewData:     newDataJSON, // Pasti berisi minimal "{}"
		IPAddress:   ipAddress,
	}

	// Fail-Safe: Jika gagal log, jangan batalkan transaksi utama
	if err := tx.Create(&activityLog).Error; err != nil {
		log.Printf("⚠️ [Activity Log Error] Gagal mencatat aktivitas di modul %s: %v\n", module, err)
	}
}