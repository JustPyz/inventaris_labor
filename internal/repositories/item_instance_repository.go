package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type ItemInstanceRepository struct {
	db *gorm.DB
}

func NewItemInstanceRepository(db *gorm.DB) *ItemInstanceRepository {
	return &ItemInstanceRepository{db: db}
}

func (r *ItemInstanceRepository) List(items *[]models.ItemInstance) error {
	return r.db.Preload("Perangkat").Order("id ASC").Find(items).Error
}

func (r *ItemInstanceRepository) Create(item *models.ItemInstance) error {
	return r.db.Create(item).Error
}

func (r *ItemInstanceRepository) FindByID(id uint, item *models.ItemInstance) error {
	return r.db.Preload("Perangkat").First(item, id).Error
}

func (r *ItemInstanceRepository) Update(item *models.ItemInstance) error {
	return r.db.Save(item).Error
}

func (r *ItemInstanceRepository) Delete(item *models.ItemInstance) error {
	return r.db.Delete(item).Error
}

func (r *ItemInstanceRepository) PerangkatExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Perangkat{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ItemInstanceRepository) KodeAssetExists(kodeAsset string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ItemInstance{}).Where("kode_asset = ?", kodeAsset).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ItemInstanceRepository) KodeAssetExistsExceptID(id uint, kodeAsset string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ItemInstance{}).
		Where("kode_asset = ? AND id <> ?", kodeAsset, id).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
