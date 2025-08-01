package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte

func InitSecret() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file tidak ditemukan, gunakan environment variables sistem")
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JwtSecret) == 0 {
		log.Fatal("❌ JWT_SECRET tidak ditemukan di .env")
	}
}
