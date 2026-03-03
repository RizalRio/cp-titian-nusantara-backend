package utils

import (
	"log"
	"os"
	"strings"
)

// DeletePhysicalFile akan menghapus file fisik di dalam direktori server
func DeletePhysicalFile(fileURL string) {
	if fileURL == "" {
		return
	}

	// Pisahkan berdasarkan path "/uploads/"
	// Contoh: http://localhost:8080/uploads/images/abc.png -> ["http://localhost:8080", "images/abc.png"]
	parts := strings.Split(fileURL, "/uploads/")
	if len(parts) < 2 {
		return // Bukan file lokal dari folder uploads kita
	}

	// Rangkai ulang path relatif ke direktori lokal
	localPath := "uploads/" + parts[1]

	// Hapus file
	err := os.Remove(localPath)
	if err != nil {
		// Log error tapi jangan hentikan eksekusi program (fail-safe)
		log.Printf("⚠️ Gagal menghapus file fisik %s: %v\n", localPath, err)
	} else {
		log.Printf("✅ File fisik berhasil dihapus: %s\n", localPath)
	}
}