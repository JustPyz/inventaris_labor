package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type KategoriRepository struct {
	db *gorm.DB
}

func NewKategoriRepository(db *gorm.DB) *KategoriRepository {
	return &KategoriRepository{db: db}
}

func (r *KategoriRepository) Create(kategori *models.Kategori) error {
	return r.db.Create(kategori).Error
}

func (r *KategoriRepository) List(kategori *[]models.Kategori) error {
	return r.db.Order("id ASC").Find(kategori).Error
}

func (r *KategoriRepository) FindByID(id uint, kategori *models.Kategori) error {
	return r.db.First(kategori, id).Error
}

func (r *KategoriRepository) Update(kategori *models.Kategori) error {
	return r.db.Save(kategori).Error
}

func (r *KategoriRepository) Delete(kategori *models.Kategori) error {
	return r.db.Delete(kategori).Error
}
