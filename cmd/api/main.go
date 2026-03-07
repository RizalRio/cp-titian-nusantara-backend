package main

import (
	"log"
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
	authService := services.NewAuthService(userRepo, config.DB)
	authHandler := handlers.NewAuthHandler(authService)

	dashboardService := services.NewDashboardService(config.DB)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// 🌟 INJEKSI DEPENDENSI BARU UNTUK PAGES
	pageRepo := repositories.NewPageRepository(config.DB)
	pageService := services.NewPageService(pageRepo, config.DB) // Pastikan PageService menerima db *gorm.DB untuk transaksi log
	pageHandler := handlers.NewPageHandler(pageService)

	// 🌟 INJEKSI DEPENDENSI UNTUK SETTINGS
	settingRepo := repositories.NewSettingRepository(config.DB)
	settingService := services.NewSettingService(settingRepo, config.DB) // Pastikan SettingService menerima db *gorm.DB untuk transaksi log
	settingHandler := handlers.NewSettingHandler(settingService)

	// 🌟 INJEKSI DEPENDENSI UNTUK CATEGORY, TAG, POST
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	tagRepo := repositories.NewTagRepository(config.DB)
	postRepo := repositories.NewPostRepository(config.DB)

	// Inisialisasi Service untuk Category, Tag, dan Post
	categoryService := services.NewCategoryService(categoryRepo, config.DB) // Pastikan CategoryService menerima db *gorm.DB untuk transaksi log
	tagService := services.NewTagService(tagRepo, config.DB)             // Pastikan TagService menerima db *gorm.DB untuk transaksi log
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

	// 🌟 INISIALISASI PROJECT (Jejak Karya)
	projectRepo := repositories.NewProjectRepository(config.DB)
	projectService := services.NewProjectService(projectRepo, config.DB)
	projectHandler := handlers.NewProjectHandler(projectService)

	// 🌟 INISIALISASI JEJAK KARYA (PORTFOLIOS) LINTAS SEKTOR
	portfolioRepo := repositories.NewPortfolioRepository(config.DB)
	portfolioService := services.NewPortfolioService(portfolioRepo, config.DB)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)

	contactRepo := repositories.NewContactRepository(config.DB)
	contactService := services.NewContactService(contactRepo, config.DB) // Pastikan ContactService menerima db *gorm.DB untuk transaksi log
	contactHandler := handlers.NewContactHandler(contactService)

	activityLogRepo := repositories.NewActivityLogRepository(config.DB)
	activityLogService := services.NewActivityLogService(activityLogRepo)
	activityLogHandler := handlers.NewActivityLogHandler(activityLogService)

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
		// 🌟 ENDPOINT AUTHENTICATION
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
		}

		// 🌟 ENDPOINT PUBLIK (Tanpa Middleware Token)
		v1.GET("/settings", settingHandler.GetSettings)
		
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

		// 🌟 ENDPOINT PROJECT
		v1.GET("/projects", projectHandler.GetAll)
		v1.GET("/projects/:id", projectHandler.GetByID)
		v1.GET("/projects/slug/:slug", projectHandler.GetBySlug)

		// 🌟 ENDPOINT PORTFOLIO (Jejak Karya)
		v1.GET("/portfolios", portfolioHandler.GetAll)
		v1.GET("/portfolios/:id", portfolioHandler.GetByID)
		v1.GET("/portfolios/slug/:slug", portfolioHandler.GetBySlug)
		
		// 🌟 ENDPOINT KONTAK
		v1.POST("/contact-messages", contactHandler.SubmitMessage)
		v1.POST("/collaboration-requests", contactHandler.SubmitCollaboration)

		// 🌟 ENDPOINT ADMIN (Dilindungi Middleware)
		adminGroup := v1.Group("/admin")
		adminGroup.Use(middleware.RequireAuth()) // Pasang satpam di sini!
		{
			adminGroup.GET("/dashboard-stats", dashboardHandler.GetStats)
			adminGroup.POST("/auth/logout", authHandler.Logout)

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

			// 🌟 CRUD PROJECT
			adminGroup.POST("/projects", projectHandler.Create)
			adminGroup.PUT("/projects/:id", projectHandler.Update)
			adminGroup.DELETE("/projects/:id", projectHandler.Delete)

			// 🌟 CRUD JEJAK KARYA
			adminGroup.POST("/portfolios", portfolioHandler.Create)
			adminGroup.PUT("/portfolios/:id", portfolioHandler.Update)
			adminGroup.DELETE("/portfolios/:id", portfolioHandler.Delete)

			// 🌟 MANAJEMEN PESAN UMUM
			adminGroup.GET("/contact-messages", contactHandler.GetAllMessages)
			adminGroup.PUT("/contact-messages/:id/read", contactHandler.MarkMessageAsRead)
			adminGroup.DELETE("/contact-messages/:id", contactHandler.DeleteMessage)

			// 🌟 MANAJEMEN KOLABORASI
			adminGroup.GET("/collaboration-requests", contactHandler.GetAllCollaborations)
			adminGroup.PUT("/collaboration-requests/:id/status", contactHandler.UpdateCollaborationStatus)
			adminGroup.DELETE("/collaboration-requests/:id", contactHandler.DeleteCollaboration)
		
			adminGroup.GET("/activity-logs", activityLogHandler.GetAllLogs)
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