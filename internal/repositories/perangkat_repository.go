package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type PerangkatRepository struct {
	db *gorm.DB
}

func NewPerangkatRepository(db *gorm.DB) *PerangkatRepository {
	return &PerangkatRepository{db: db}
}

func (r *PerangkatRepository) List(perangkat *[]models.Perangkat) error {
	return r.db.Order("id ASC").Find(perangkat).Error
}

func (r *PerangkatRepository) Create(perangkat *models.Perangkat) error {
	return r.db.Create(perangkat).Error
}

func (r *PerangkatRepository) FindByID(id uint, perangkat *models.Perangkat) error {
	return r.db.First(perangkat, id).Error
}

func (r *PerangkatRepository) Update(perangkat *models.Perangkat) error {
	return r.db.Save(perangkat).Error
}

func (r *PerangkatRepository) Delete(perangkat *models.Perangkat) error {
	return r.db.Delete(perangkat).Error
}

func (r *PerangkatRepository) KategoriExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Kategori{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PerangkatRepository) JurusanExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Jurusan{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PerangkatRepository) LaborExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Labor{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
