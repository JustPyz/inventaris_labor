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
	GetStatsByItemInstanceID(itemInstanceID uint) (*models.KerusakanStats, error)
	// HasOpenKerusakan mengecek apakah item_instance masih punya kerusakan
	// yang belum berstatus 'selesai', di luar kerusakan dengan ID yang dikecualikan.
	HasOpenKerusakan(itemInstanceID uint, excludeKerusakanID uint) (bool, error)
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

// HasOpenKerusakan mengembalikan true jika masih ada kerusakan dengan status
// bukan 'selesai' pada item_instance yang dimaksud (di luar kerusakanID yang dikecualikan).
func (r *kerusakanRepository) HasOpenKerusakan(itemInstanceID uint, excludeKerusakanID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Kerusakan{}).
		Where("id_item_instance = ? AND status != ? AND id != ?", itemInstanceID, "selesai", excludeKerusakanID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *kerusakanRepository) ItemInstanceExists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ItemInstance{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *kerusakanRepository) GetStatsByItemInstanceID(itemInstanceID uint) (*models.KerusakanStats, error) {
	var totalPerbaikan int64
	var totalBiaya int64

	row := r.db.Table("perbaikans").
		Joins("JOIN kerusakans ON perbaikans.id_kerusakan = kerusakans.id").
		Where("kerusakans.id_item_instance = ?", itemInstanceID).
		Select("count(perbaikans.id), coalesce(sum(perbaikans.biaya), 0)").Row()

	err := row.Scan(&totalPerbaikan, &totalBiaya)
	if err != nil {
		return nil, err
	}

	return &models.KerusakanStats{
		ItemInstanceID: itemInstanceID,
		TotalPerbaikan: totalPerbaikan,
		TotalBiaya:     totalBiaya,
	}, nil
}
