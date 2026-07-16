package repositories

import (
	"invela-be/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) List(roles *[]models.Role) error {
	return r.db.Order("id ASC").Find(roles).Error
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) FindByID(id uint, role *models.Role) error {
	return r.db.First(role, id).Error
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(role *models.Role) error {
	return r.db.Delete(role).Error
}
