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

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// 🌟 BARU: Mengambil semua user (tanpa password hash)
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Select("id, name, email, role_id, status, last_login_at, created_at").Find(&users).Error
	return users, err
}

// 🌟 BARU: Mengambil detail satu user
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}