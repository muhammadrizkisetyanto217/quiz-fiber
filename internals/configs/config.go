// package configs

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// // LoadEnv memuat file .env
// func LoadEnv() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("❌ Gagal memuat .env file. Pastikan file .env tersedia!")
// 	}
// 	log.Println("✅ .env file berhasil dimuat!")
// }

// // GetEnv mengambil nilai dari .env dengan default value
// func GetEnv(key string, defaultValue ...string) string {
// 	value, exists := os.LookupEnv(key)
// 	if !exists && len(defaultValue) > 0 {
// 		return defaultValue[0]
// 	}
// 	return value
// }

package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// JWTSecret akan diisi saat LoadEnv dipanggil
var JWTSecret string

// LoadEnv memuat file .env jika berjalan secara lokal
func LoadEnv() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" { // Cek apakah berjalan di lokal
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Tidak menemukan .env file, menggunakan environment variable dari sistem")
		} else {
			log.Println("✅ .env file berhasil dimuat!")
		}
	} else {
		log.Println("🚀 Running in Railway, menggunakan environment variables dari sistem")
	}

	// 🔐 Ambil JWT_SECRET setelah dotenv diload
	JWTSecret = GetEnv("JWT_SECRET")

	// Debug print untuk memastikan terload
	if JWTSecret == "" {
		log.Println("❌ JWT_SECRET belum diset! Harap cek .env atau Environment Variable di Railway.")
	} else {
		log.Println("✅ JWT_SECRET berhasil dimuat.")
	}
}

// GetEnv mengambil nilai dari .env dengan default value
func GetEnv(key string, defaultValue ...string) string {
	value, exists := os.LookupEnv(key)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}
