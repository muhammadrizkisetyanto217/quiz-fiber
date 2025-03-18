package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"quiz-fiber/internals/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB untuk menghubungkan ke database
func ConnectDB() error {
	// ✅ Load environment variables
	configs.LoadEnv()

	// ✅ Ambil URL dari environment variable
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return fmt.Errorf("❌ Database URL is not set in environment variables")
	}

	// ✅ Konfigurasi koneksi database dengan GORM (Logging Level: Info)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 🔥 Logging query yang dieksekusi
	})
	if err != nil {
		return fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	// ✅ Setup Connection Pooling
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("❌ Failed to get database instance: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)                 // Maksimal 20 koneksi
	sqlDB.SetMaxIdleConns(10)                 // Maksimal 10 koneksi idle
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Waktu maksimal koneksi hidup
	sqlDB.SetConnMaxIdleTime(3 * time.Minute) // Koneksi idle lebih dari 3 menit akan ditutup

	// ✅ Logging untuk memastikan koneksi berhasil
	log.Println("✅ Database connected successfully!")

	DB = db

	// 🔥 Coba jalankan query SELECT sederhana untuk debug
	// var columns []string
	// result := DB.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'themes_or_levels'").Scan(&columns)
	// if result.Error != nil {
	// 	log.Println("❌ Error fetching column names:", result.Error)
	// } else {
	// 	log.Println("✅ Columns in themes_or_levels:", columns)
	// }

	return nil
}
