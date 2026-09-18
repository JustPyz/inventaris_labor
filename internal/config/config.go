package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// defaultDevJWTSecret hanya boleh dipakai di lingkungan non-produksi.
// Di produksi (APP_ENV=production) secret ini ditolak dan JWT_SECRET wajib di-set.
const defaultDevJWTSecret = "invela-secret-key-change-in-production"

type Config struct {
	Port          string
	DBPath        string
	JWTSecret     string
	Env           string
	AdminUsername string
	AdminPassword string
}

func Load() Config {
	// Load .env file jika ada (tidak error jika tidak ditemukan)
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: .env file not found, using system environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/app.db"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Di produksi, JWT_SECRET wajib. Tanpa ini, token bisa ditempa
		// menggunakan secret default yang publik di repo.
		if env == "production" {
			log.Fatal("JWT_SECRET environment variable is required when APP_ENV=production")
		}
		log.Printf("WARNING: JWT_SECRET not set, using insecure default secret (APP_ENV=%s). DO NOT use this in production.", env)
		jwtSecret = defaultDevJWTSecret
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin@smkn4pyk.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	return Config{
		Port:          port,
		DBPath:        dbPath,
		JWTSecret:     jwtSecret,
		Env:           env,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
	}
}
