package repositories

import (
	"invela-be/internal/models"
	"gorm.io/gorm"
)

type PerbaikanRepository interface {
	Create(perbaikan *models.Perbaikan) error
	FindAll() ([]models.Perbaikan, error)
	FindByID(id uint) (*models.Perbaikan, error)
	Update(perbaikan *models.Perbaikan) error
	Delete(id uint) error
	KerusakanExists(id uint) (bool, error)
}

type perbaikanRepository struct {
	db *gorm.DB
}

func NewPerbaikanRepository(db *gorm.DB) PerbaikanRepository {
	return &perbaikanRepository{db}
}

func (r *perbaikanRepository) Create(perbaikan *models.Perbaikan) error {
	return r.db.Create(perbaikan).Error
}

func (r *perbaikanRepository) FindAll() ([]models.Perbaikan, error) {
	var perbaikans []models.Perbaikan
	err := r.db.Preload("User").Preload("Kerusakan").Preload("Kerusakan.ItemInstance").Preload("Kerusakan.ItemInstance.Perangkat").Find(&perbaikans).Error
	return perbaikans, err
}

func (r *perbaikanRepository) FindByID(id uint) (*models.Perbaikan, error) {
	var perbaikan models.Perbaikan
	err := r.db.Preload("User").Preload("Kerusakan").Preload("Kerusakan.ItemInstance").Preload("Kerusakan.ItemInstance.Perangkat").First(&perbaikan, id).Error
	if err != nil {
		return nil, err
	}
	return &perbaikan, nil
}

func (r *perbaikanRepository) Update(perbaikan *models.Perbaikan) error {
	return r.db.Save(perbaikan).Error
}

func (r *perbaikanRepository) Delete(id uint) error {
	return r.db.Delete(&models.Perbaikan{}, id).Error
}

func (r *perbaikanRepository) KerusakanExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Kerusakan{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
