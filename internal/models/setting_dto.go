package models

// DTO: Kita gunakan array of struct agar Admin bisa update banyak pengaturan sekaligus (Batch Update)
type UpsertSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"` // Jika kosong, setidaknya string kosong ""
}