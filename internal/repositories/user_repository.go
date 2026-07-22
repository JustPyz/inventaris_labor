package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) List(users *[]models.User) error {
	return r.db.Preload("Role").Preload("Jurusan").Order("id ASC").Find(users).Error
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint, user *models.User) error {
	return r.db.Preload("Role").Preload("Jurusan").First(user, id).Error
}

func (r *UserRepository) RoleExists(roleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Role{}).Where("id = ?", roleID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) JurusanExists(jurusanID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Jurusan{}).Where("id = ?", jurusanID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *UserRepository) FindByUsername(username string, user *models.User) error {
	return r.db.Preload("Role").Where("username = ?", username).First(user).Error
}
