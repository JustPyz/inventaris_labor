package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type PenggunaanRepository struct {
	db *gorm.DB
}

func NewPenggunaanRepository(db *gorm.DB) *PenggunaanRepository {
	return &PenggunaanRepository{db: db}
}

func (r *PenggunaanRepository) List(items *[]models.Penggunaan) error {
	return r.db.Preload("Labor").Preload("Kelas").Order("id ASC").Find(items).Error
}

func (r *PenggunaanRepository) Create(item *models.Penggunaan) error {
	return r.db.Create(item).Error
}

func (r *PenggunaanRepository) FindByID(id uint, item *models.Penggunaan) error {
	return r.db.Preload("Labor").Preload("Kelas").First(item, id).Error
}

func (r *PenggunaanRepository) Delete(item *models.Penggunaan) error {
	return r.db.Delete(item).Error
}

func (r *PenggunaanRepository) UserExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PenggunaanRepository) FindUserByID(id uint, user *models.User) error {
	return r.db.First(user, id).Error
}

func (r *PenggunaanRepository) LaborExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Labor{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PenggunaanRepository) KelasExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Kelas{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
