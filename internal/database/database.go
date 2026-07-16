package database

import (
	"fmt"
	"os"
	"path/filepath"

	"invela-be/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Open(dbPath string) (*gorm.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.AutoMigrate(&models.Role{}); err != nil {
		return nil, fmt.Errorf("auto migrate role: %w", err)
	}

	return db, nil
}
