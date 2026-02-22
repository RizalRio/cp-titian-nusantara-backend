package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/jwt"
	"backend/pkg/utils"
	"errors"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

func (s *AuthService) Login(req models.LoginRequest) (string, *models.User, error) {
	// 1. Cari user di database
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// PENTING: Jangan beri tahu detail "Email tidak ditemukan". Gunakan pesan generik.
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

	return token, user, nil
}