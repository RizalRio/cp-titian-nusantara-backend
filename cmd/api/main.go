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

	"time"

	"github.com/gin-contrib/cors"
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

	// 🌟 INJEKSI DEPENDENSI UNTUK SETTINGS
	settingRepo := repositories.NewSettingRepository(config.DB)
	settingService := services.NewSettingService(settingRepo)
	settingHandler := handlers.NewSettingHandler(settingService)

	// 🌟 INJEKSI DEPENDENSI UNTUK CATEGORY, TAG, POST
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	tagRepo := repositories.NewTagRepository(config.DB)
	postRepo := repositories.NewPostRepository(config.DB)

	// Inisialisasi Service untuk Category, Tag, dan Post
	categoryService := services.NewCategoryService(categoryRepo)
	tagService := services.NewTagService(tagRepo)
	postService := services.NewPostService(postRepo, config.DB) // PostService butuh koneksi DB langsung

	// Inisialisasi Handler untuk Category, Tag, dan Post
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	tagHandler := handlers.NewTagHandler(tagService)
	postHandler := handlers.NewPostHandler(postService)

	mediaHandler := handlers.NewMediaHandler()

	// 🌟 INISIALISASI EKOSISTEM LAYANAN
	serviceRepo := repositories.NewServiceRepository(config.DB)
	serviceEcosystemService := services.NewServiceEcosystemService(serviceRepo, config.DB)
	serviceHandler := handlers.NewServiceHandler(serviceEcosystemService)

	// 3. Setup Framework Gin (Router)
	r := gin.Default()

	// 🌟 EKSPOS FOLDER STATIS
	// Akses file lokal di ./uploads akan diarahkan ke route /uploads
	r.Static("/uploads", "./uploads")

	// 🌟 TAMBAHKAN MIDDLEWARE CORS DI SINI 🌟
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Izinkan Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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

		// 🌟 ENDPOINT PUBLIK (Tanpa Middleware Token)
		v1.GET("/settings", settingHandler.GetPublicSettings)
		
		// 👇 TAMBAHKAN RUTE INI DI SINI 👇
		v1.GET("/pages/:slug", pageHandler.GetBySlug)

		// 🌍 RUTING PUBLIK (Tanpa Token - Untuk dibaca oleh Pengunjung Website)
		v1.GET("/categories", categoryHandler.GetAll)
		v1.GET("/tags", tagHandler.GetAll)
		v1.GET("/posts", postHandler.GetAll)      // Bisa ditambah query: ?status=published
		v1.GET("/posts/:id", postHandler.GetByID) // Untuk admin/editor
		v1.GET("/posts/slug/:slug", postHandler.GetBySlug) // Untuk halaman detail pembaca

		// 🌟 ENDPOINT LAYANAN (Bisa diakses publik untuk melihat daftar layanan, tapi hanya admin yang bisa membuat/mengedit)
		v1.GET("/services", serviceHandler.GetAll)
		v1.GET("/services/:id", serviceHandler.GetByID)
		v1.GET("/services/slug/:slug", serviceHandler.GetBySlug)

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

			// 🌟 ROUTE UNTUK UPDATE SETTINGS (Batch Update)
			adminGroup.PUT("/settings", settingHandler.UpdateSettings)

			// Endpoint Kategori & Tag Admin
			adminGroup.POST("/categories", categoryHandler.Create)
			adminGroup.PUT("/categories/:id", categoryHandler.Update)
			adminGroup.DELETE("/categories/:id", categoryHandler.Delete)

			adminGroup.POST("/tags", tagHandler.Create)
			adminGroup.PUT("/tags/:id", tagHandler.Update)
			adminGroup.DELETE("/tags/:id", tagHandler.Delete)

			// Endpoint Post Admin
			adminGroup.POST("/posts", postHandler.Create)
			adminGroup.PUT("/posts/:id", postHandler.Update)
			adminGroup.DELETE("/posts/:id", postHandler.Delete)

			adminGroup.POST("/media/upload", mediaHandler.UploadImage)

			// 🌟 CRUD EKOSISTEM LAYANAN
			adminGroup.POST("/services", serviceHandler.Create)
			adminGroup.PUT("/services/:id", serviceHandler.Update)
			adminGroup.DELETE("/services/:id", serviceHandler.Delete)
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