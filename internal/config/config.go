package config

import "os"

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/app.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "invela-secret-key-change-in-production"
	}

	return Config{
		Port:      port,
		DBPath:    dbPath,
		JWTSecret: jwtSecret,
	}
}
