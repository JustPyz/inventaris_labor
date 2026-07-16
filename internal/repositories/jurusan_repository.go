package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type JurusanRepository struct {
	db *gorm.DB
}

func NewJurusanRepository(db *gorm.DB) *JurusanRepository {
	return &JurusanRepository{db: db}
}

func (r *JurusanRepository) List(jurusan *[]models.Jurusan) error {
	return r.db.Order("id ASC").Find(jurusan).Error
}

func (r *JurusanRepository) Create(jurusan *models.Jurusan) error {
	return r.db.Create(jurusan).Error
}

func (r *JurusanRepository) FindByID(id uint, jurusan *models.Jurusan) error {
	return r.db.First(jurusan, id).Error
}

func (r *JurusanRepository) Update(jurusan *models.Jurusan) error {
	return r.db.Save(jurusan).Error
}

func (r *JurusanRepository) Delete(jurusan *models.Jurusan) error {
	return r.db.Delete(jurusan).Error
}
