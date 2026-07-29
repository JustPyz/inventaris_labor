package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type RiwayatPerbaikanRepository interface {
	Create(r *models.RiwayatPerbaikan) error
	FindAll() ([]models.RiwayatPerbaikan, error)
	FindByKodeAsset(kodeAsset string) ([]models.RiwayatPerbaikan, error)
}

type riwayatPerbaikanRepository struct {
	db *gorm.DB
}

func NewRiwayatPerbaikanRepository(db *gorm.DB) RiwayatPerbaikanRepository {
	return &riwayatPerbaikanRepository{db}
}

func (r *riwayatPerbaikanRepository) Create(rw *models.RiwayatPerbaikan) error {
	return r.db.Create(rw).Error
}

func (r *riwayatPerbaikanRepository) FindAll() ([]models.RiwayatPerbaikan, error) {
	var results []models.RiwayatPerbaikan
	err := r.db.Order("tanggal_perbaikan DESC").Find(&results).Error
	return results, err
}

func (r *riwayatPerbaikanRepository) FindByKodeAsset(kodeAsset string) ([]models.RiwayatPerbaikan, error) {
	var results []models.RiwayatPerbaikan
	err := r.db.Where("kode_asset = ?", kodeAsset).
		Order("tanggal_perbaikan DESC").
		Find(&results).Error
	return results, err
}
