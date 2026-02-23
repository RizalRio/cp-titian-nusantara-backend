package main

import (
	"log"
	"net/http"
	"os"

	"backend/config" // Sesuaikan dengan nama module di go.mod kamu

	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repositories"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ File .env tidak ditemukan, menggunakan environment variable bawaan sistem")
	}

	// 2. Inisialisasi Koneksi Database
	config.ConnectDB()

	userRepo := repositories.NewUserRepository(config.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// 🌟 INJEKSI DEPENDENSI BARU UNTUK PAGES
	pageRepo := repositories.NewPageRepository(config.DB)
	pageService := services.NewPageService(pageRepo)
	pageHandler := handlers.NewPageHandler(pageService)

	// 3. Setup Framework Gin (Router)
	r := gin.Default()

	// 4. Grouping Route untuk API versi 1
	v1 := r.Group("/api/v1")
	{
		// 🌟 ENDPOINT HEALTH CHECK
		// Berfungsi untuk mengecek apakah server menyala dan database merespons
		v1.GET("/health", func(c *gin.Context) {
			// Cek "denyut nadi" (ping) fisik ke database
			sqlDB, err := config.DB.DB()
			dbStatus := "connected"
			
			if err != nil || sqlDB.Ping() != nil {
				dbStatus = "disconnected"
			}

			// Kembalikan response JSON
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"message": "Sistem Titian Nusantara berjalan dengan baik 🚀",
				"database": dbStatus,
			})
		})

		// 🌟 ENDPOINT AUTHENTICATION
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
		}

		// 🌟 ENDPOINT ADMIN (Dilindungi Middleware)
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.RequireAuth()) // Pasang satpam di sini!
		{
			// Contoh rute untuk mengecek profil admin yang sedang login
			adminGroup.GET("/me", func(c *gin.Context) {
				// Ambil user_id dari hasil ekstrak token di middleware
				userID := c.MustGet("user_id").(string)

				c.JSON(http.StatusOK, gin.H{
					"status":  "success",
					"message": "Akses Admin diizinkan!",
					"data": gin.H{
						"user_id": userID,
					},
				})
			})

			// 🌟 CRUD ROUTES UNTUK PAGES
			pagesGroup := adminGroup.Group("/pages")
			{
				pagesGroup.POST("", pageHandler.Create)
				pagesGroup.GET("", pageHandler.GetAll)
				pagesGroup.GET("/:id", pageHandler.GetByID)
				pagesGroup.PUT("/:id", pageHandler.Update)
				pagesGroup.DELETE("/:id", pageHandler.Delete)
			}
		}
	}

	// 5. Menentukan Port Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default fallback jika di .env kosong
	}

	// 6. Jalankan Server
	log.Printf("🚀 Server siap menerima request di http://localhost:%s\n", port)
	r.Run(":" + port)
}