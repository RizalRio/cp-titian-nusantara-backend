package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengenkripsi password plain text menjadi hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Cost 14 sudah cukup kuat
	return string(bytes), err
}

// CheckPasswordHash membandingkan password input dengan hash di database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Jika nil, berarti password cocok
}