package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah variabel global agar koneksi database bisa dipakai di file/layer lain
var DB *gorm.DB

func ConnectDB() {
	// 1. Mengambil data dari file .env
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// 2. Merakit Data Source Name (DSN) format PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// 3. Membuka koneksi menggunakan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal terkoneksi ke Database: ", err)
	}

	log.Println("✅ Database PostgreSQL berhasil terhubung!")

	// err = database.AutoMigrate(
	// 	// &models.Page{}, // Nanti kita aktifkan saat mengerjakan halaman
	// 	&models.User{},
	// )
	// if err != nil {
	// 	log.Fatal("❌ Gagal migrasi database: ", err)
	// }
	
	// Simpan instance koneksi ke variabel global
	DB = database
}