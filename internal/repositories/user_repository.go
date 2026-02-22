package repositories

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// FindByEmail mencari user berdasarkan email di database
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	// First() akan mengembalikan error jika data tidak ditemukan (RecordNotFound)
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}