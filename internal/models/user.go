package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex:idx_username_deleted,composite:username_deleted" json:"username"`
	PasswordHash string         `gorm:"column:password_hash;size:255;not null" json:"-"`
	RoleID       uint           `gorm:"column:role_id;not null" json:"-"`
	Role         Role           `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	JurusanID    *uint          `gorm:"column:jurusan_id" json:"jurusan_id"`
	Jurusan      *Jurusan       `gorm:"foreignKey:JurusanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"jurusan"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"uniqueIndex:idx_username_deleted,composite:username_deleted" json:"-"`
}
