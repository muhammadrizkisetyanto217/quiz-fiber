package main

import (
	"log"

	"quiz-fiber/internals/configs"
	"quiz-fiber/internals/seeds"
)

func main() {
	configs.LoadEnv() // <-- penting
	db := configs.InitDB()
	log.Println("🚀 Menjalankan semua seed...")
	seeds.RunAllSeeds(db)
}
