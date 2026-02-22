package middleware

import (
	"net/http"
	"strings"

	"backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RequireAuth adalah satpam untuk rute yang membutuhkan login
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")

		// 2. Cek apakah header kosong
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Akses ditolak. Token tidak ditemukan.",
			})
			return
		}

		// 3. Pastikan formatnya "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Format token tidak valid. Gunakan format: Bearer <token>",
			})
			return
		}

		// 4. Validasi Token
		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Sesi telah berakhir atau token tidak valid. Silakan login kembali.",
			})
			return
		}

		// 5. Simpan data user ke dalam context (memori sementara untuk request ini)
		// Jadi di handler API nanti, kita bisa langsung panggil c.MustGet("user_id")
		c.Set("user_id", claims["user_id"])
		c.Set("role_id", claims["role_id"])

		// 6. Lanjutkan ke rute tujuan (silakan masuk)
		c.Next()
	}
}