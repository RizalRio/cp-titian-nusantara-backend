package models

// DTO untuk pembuatan User baru
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id"` // Kosongkan jika belum ada sistem Role yang kompleks
	Status   string `json:"status" binding:"required,oneof=active inactive suspended"`
}

// DTO untuk pembaruan User
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"` // Opsional, hanya diisi jika ingin ganti password
	RoleID   string `json:"role_id"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}