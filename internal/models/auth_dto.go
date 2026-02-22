package models

// LoginRequest adalah format JSON yang diharapkan dari Frontend Next.js
type LoginRequest struct {
	// binding:"required,email" otomatis memvalidasi apakah input wajib diisi dan berformat email
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}