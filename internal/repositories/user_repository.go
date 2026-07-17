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
	return r.db.Preload("Role").Order("id ASC").Find(users).Error
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint, user *models.User) error {
	return r.db.Preload("Role").First(user, id).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}
