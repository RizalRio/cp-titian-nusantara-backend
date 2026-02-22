package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken membuat token JWT dengan masa berlaku tertentu
func GenerateToken(userID string, roleID string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	// Payload token (isi data di dalam token)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role_id": roleID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expired dalam 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key dari .env
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}