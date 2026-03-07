package services

import (
	"errors"
	"strings"

	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/utils" // Pastikan utilitas HashPassword Anda ada di sini

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
	db   *gorm.DB
}

func NewUserService(repo *repositories.UserRepository, db *gorm.DB) *UserService {
	return &UserService{repo: repo, db: db}
}

func (s *UserService) CreateUser(req models.CreateUserRequest, actorID *uuid.UUID, ipAddress string) (*models.User, error) {
	// Hash password sebelum disimpan
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses kata sandi")
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       req.Status,
	}

	if req.RoleID != "" {
		user.RoleID = &req.RoleID
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// 🌟 CATAT LOG AKTIVITAS
		LogActivity(tx, actorID, "CREATE", "Users", "Menambahkan pengguna baru: "+user.Email, ipAddress, nil, user)
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("email sudah terdaftar, gunakan email lain")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest, actorID *uuid.UUID, ipAddress string) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	oldDataSnapshot := *user

	if req.Name != "" { user.Name = req.Name }
	if req.Email != "" { user.Email = req.Email }
	if req.Status != "" { user.Status = req.Status }
	if req.RoleID != "" { user.RoleID = &req.RoleID }

	// Hanya hash dan update password jika diisi oleh Admin
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("gagal memproses pembaruan kata sandi")
		}
		user.PasswordHash = hashedPassword
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		// 🌟 CATAT LOG AKTIVITAS
		LogActivity(tx, actorID, "UPDATE", "Users", "Memperbarui data pengguna: "+user.Email, ipAddress, oldDataSnapshot, user)
		return nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("email sudah digunakan oleh pengguna lain")
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id string, actorID *uuid.UUID, ipAddress string) error {
	userToDelete, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	// Mencegah admin menghapus dirinya sendiri (Fail-safe)
	if actorID != nil && id == actorID.String() {
		return errors.New("anda tidak dapat menghapus akun Anda sendiri saat sedang login")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", id).Error; err != nil {
			return err
		}
		// 🌟 CATAT LOG AKTIVITAS
		LogActivity(tx, actorID, "DELETE", "Users", "Menghapus pengguna: "+userToDelete.Email, ipAddress, userToDelete, nil)
		return nil
	})
}