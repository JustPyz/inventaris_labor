package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type LaborRepository struct {
	db *gorm.DB
}

func NewLaborRepository(db *gorm.DB) *LaborRepository {
	return &LaborRepository{db: db}
}

func (r *LaborRepository) List(labor *[]models.Labor) error {
	return r.db.Order("id ASC").Find(labor).Error
}

func (r *LaborRepository) Create(labor *models.Labor) error {
	return r.db.Create(labor).Error
}

func (r *LaborRepository) FindByID(id uint, labor *models.Labor) error {
	return r.db.First(labor, id).Error
}

func (r *LaborRepository) Update(labor *models.Labor) error {
	return r.db.Save(labor).Error
}

func (r *LaborRepository) Delete(labor *models.Labor) error {
	return r.db.Delete(labor).Error
}
