package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/jwt"
	"backend/pkg/utils"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	db       *gorm.DB // 🌟 INJEKSI: Diperlukan untuk log aktivitas
}

// 🌟 INJEKSI: Tambahkan db *gorm.DB
func NewAuthService(repo *repositories.UserRepository, db *gorm.DB) *AuthService {
	return &AuthService{userRepo: repo, db: db}
}

// 🌟 INJEKSI LOG: Tambahkan parameter ipAddress
func (s *AuthService) Login(req models.LoginRequest, ipAddress string) (string, *models.User, error) {
	// 1. Cari user di database
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("Email atau password salah")
	}

	// 2. Bandingkan password plain text dari frontend dengan Hash di database
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", nil, errors.New("Email atau password salah")
	}

	// 3. Pastikan akun tidak di-suspend
	if user.Status != "active" {
		return "", nil, errors.New("Akun Anda tidak aktif, silakan hubungi Administrator")
	}

	// 4. Generate JWT Token
	roleID := ""
	if user.RoleID != nil {
		roleID = *user.RoleID
	}
	
	token, err := jwt.GenerateToken(user.ID, roleID)
	if err != nil {
		return "", nil, errors.New("Gagal memproses sesi login")
	}

	// 🌟 INJEKSI LOG AKTIVITAS (LOGIN)
	// Kita siapkan pointer UUID yang kosong
	var userIDPtr *uuid.UUID

	// Jika user.ID di model Anda bertipe 'string', kita parsing dulu:
	if idStr, ok := any(user.ID).(string); ok {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			userIDPtr = &parsedID
		}
	} else if idUUID, ok := any(user.ID).(uuid.UUID); ok {
		// Jika user.ID memang sudah bertipe 'uuid.UUID', langsung ambil alamat memorinya:
		userIDPtr = &idUUID
	}

	// Catat log HANYA JIKA login berhasil
	_ = s.db.Transaction(func(tx *gorm.DB) error {
		LogActivity(tx, userIDPtr, "LOGIN", "Auth", "Pengguna berhasil masuk ke dalam sistem CMS", ipAddress, nil, nil)
		return nil
	})

	return token, user, nil
}