package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Gagal memuat .env file. Pastikan file .env tersedia!")
	}
	log.Println("✅ .env file berhasil dimuat!")
}

// GetEnv mengambil nilai dari .env dengan default value
func GetEnv(key string, defaultValue ...string) string {
	value, exists := os.LookupEnv(key)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

