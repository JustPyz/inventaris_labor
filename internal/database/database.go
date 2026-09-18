package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"invela-be/internal/config"
	"invela-be/internal/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	dir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := enableForeignKeys(db); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Role{}); err != nil {
		return nil, fmt.Errorf("auto migrate role: %w", err)
	}

	if err := db.AutoMigrate(&models.Jurusan{}); err != nil {
		return nil, fmt.Errorf("auto migrate jurusan: %w", err)
	}

	// Drop the old unique index on username before migrating User.
	// The new model uses a composite unique index on (username, deleted_at)
	// to allow re-creating soft-deleted usernames.
	db.Exec("DROP INDEX IF EXISTS idx_users_username")

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("auto migrate user: %w", err)
	}

	if err := db.AutoMigrate(&models.Kelas{}); err != nil {
		return nil, fmt.Errorf("auto migrate kelas: %w", err)
	}

	if err := db.AutoMigrate(&models.Kategori{}); err != nil {
		return nil, fmt.Errorf("auto migrate kategori: %w", err)
	}

	if err := db.AutoMigrate(&models.Labor{}); err != nil {
		return nil, fmt.Errorf("auto migrate labor: %w", err)
	}

	if err := db.AutoMigrate(&models.Perangkat{}); err != nil {
		return nil, fmt.Errorf("auto migrate perangkat: %w", err)
	}

	if err := db.AutoMigrate(&models.ItemInstance{}); err != nil {
		return nil, fmt.Errorf("auto migrate item instance: %w", err)
	}

	if err := db.AutoMigrate(&models.Peminjaman{}); err != nil {
		return nil, fmt.Errorf("auto migrate peminjaman: %w", err)
	}

	if err := db.AutoMigrate(&models.Penggunaan{}); err != nil {
		return nil, fmt.Errorf("auto migrate penggunaan: %w", err)
	}

	if err := db.AutoMigrate(&models.Kerusakan{}); err != nil {
		return nil, fmt.Errorf("auto migrate kerusakan: %w", err)
	}

	if err := db.AutoMigrate(&models.Perbaikan{}); err != nil {
		return nil, fmt.Errorf("auto migrate perbaikan: %w", err)
	}

	if err := db.AutoMigrate(&models.RiwayatPerbaikan{}); err != nil {
		return nil, fmt.Errorf("auto migrate riwayat perbaikan: %w", err)
	}

	if err := seedRoles(db); err != nil {
		return nil, fmt.Errorf("seed roles: %w", err)
	}

	if err := seedAdminUser(db, cfg); err != nil {
		return nil, fmt.Errorf("seed admin user: %w", err)
	}

	return db, nil
}

func seedRoles(db *gorm.DB) error {
	roles := []string{"kabeng", "guru", "kaprog", "sapras", "admin"}

	for _, roleName := range roles {
		var existing models.Role
		result := db.Where("role = ?", roleName).First(&existing)
		if result.Error == nil {
			continue // sudah ada
		}
		if result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		if err := db.Create(&models.Role{Role: roleName}).Error; err != nil {
			return fmt.Errorf("create role %s: %w", roleName, err)
		}
		log.Printf("Seeded role: %s", roleName)
	}

	return nil
}

func seedAdminUser(db *gorm.DB, cfg config.Config) error {
	var count int64
	db.Model(&models.User{}).Where("username = ?", cfg.AdminUsername).Count(&count)
	if count > 0 {
		return nil // sudah ada
	}

	var adminRole models.Role
	if err := db.Where("role = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("find admin role: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	admin := &models.User{
		Username:     cfg.AdminUsername,
		PasswordHash: string(hash),
		RoleID:       adminRole.ID,
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	log.Printf("Seeded admin user (username: %s)\n", cfg.AdminUsername)
	return nil
}

func enableForeignKeys(db *gorm.DB) error {
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	var enabled int
	if err := db.Raw("PRAGMA foreign_keys").Scan(&enabled).Error; err != nil {
		return fmt.Errorf("verify foreign keys pragma: %w", err)
	}

	if enabled != 1 {
		return fmt.Errorf("foreign keys pragma is disabled")
	}

	return nil
}
