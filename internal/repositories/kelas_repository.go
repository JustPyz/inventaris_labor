package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type KelasRepository struct {
	db *gorm.DB
}

func NewKelasRepository(db *gorm.DB) *KelasRepository {
	return &KelasRepository{db: db}
}

func (r *KelasRepository) List(kelas *[]models.Kelas) error {
	return r.db.Order("id ASC").Find(kelas).Error
}

func (r *KelasRepository) Create(kelas *models.Kelas) error {
	return r.db.Create(kelas).Error
}

func (r *KelasRepository) FindByID(id uint, kelas *models.Kelas) error {
	return r.db.First(kelas, id).Error
}

func (r *KelasRepository) Update(kelas *models.Kelas) error {
	return r.db.Save(kelas).Error
}

func (r *KelasRepository) Delete(kelas *models.Kelas) error {
	return r.db.Delete(kelas).Error
}
