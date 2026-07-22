package repositories

import (
	"invela-be/internal/models"
	"gorm.io/gorm"
)

type KerusakanRepository interface {
	Create(kerusakan *models.Kerusakan) error
	FindAll() ([]models.Kerusakan, error)
	FindByID(id uint) (*models.Kerusakan, error)
	Update(kerusakan *models.Kerusakan) error
	Delete(id uint) error
	ItemInstanceExists(id uint) (bool, error)
}

type kerusakanRepository struct {
	db *gorm.DB
}

func NewKerusakanRepository(db *gorm.DB) KerusakanRepository {
	return &kerusakanRepository{db}
}

func (r *kerusakanRepository) Create(kerusakan *models.Kerusakan) error {
	return r.db.Create(kerusakan).Error
}

func (r *kerusakanRepository) FindAll() ([]models.Kerusakan, error) {
	var kerusakans []models.Kerusakan
	err := r.db.Preload("User").Preload("ItemInstance").Preload("ItemInstance.Perangkat").Find(&kerusakans).Error
	return kerusakans, err
}

func (r *kerusakanRepository) FindByID(id uint) (*models.Kerusakan, error) {
	var kerusakan models.Kerusakan
	err := r.db.Preload("User").Preload("ItemInstance").Preload("ItemInstance.Perangkat").First(&kerusakan, id).Error
	if err != nil {
		return nil, err
	}
	return &kerusakan, nil
}

func (r *kerusakanRepository) Update(kerusakan *models.Kerusakan) error {
	return r.db.Save(kerusakan).Error
}

func (r *kerusakanRepository) Delete(id uint) error {
	return r.db.Delete(&models.Kerusakan{}, id).Error
}

func (r *kerusakanRepository) ItemInstanceExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ItemInstance{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
