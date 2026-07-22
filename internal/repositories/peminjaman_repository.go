package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type PeminjamanRepository struct {
	db *gorm.DB
}

func NewPeminjamanRepository(db *gorm.DB) *PeminjamanRepository {
	return &PeminjamanRepository{db: db}
}

func (r *PeminjamanRepository) List(items *[]models.Peminjaman) error {
	return r.db.Preload("ItemInstance").Preload("ItemInstance.Perangkat").Order("id ASC").Find(items).Error
}

func (r *PeminjamanRepository) Create(item *models.Peminjaman) error {
	return r.db.Create(item).Error
}

func (r *PeminjamanRepository) FindByID(id uint, item *models.Peminjaman) error {
	return r.db.Preload("ItemInstance").Preload("ItemInstance.Perangkat").First(item, id).Error
}

func (r *PeminjamanRepository) Update(item *models.Peminjaman) error {
	return r.db.Save(item).Error
}

func (r *PeminjamanRepository) Delete(item *models.Peminjaman) error {
	return r.db.Delete(item).Error
}

func (r *PeminjamanRepository) ItemInstanceExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ItemInstance{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
