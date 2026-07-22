package models

import "time"

type Kerusakan struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ItemInstanceID uint         `gorm:"column:id_item_instance;not null" json:"id_item_instance"`
	ItemInstance   ItemInstance `gorm:"foreignKey:ItemInstanceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"item_instance"`
	UserID         uint         `gorm:"column:id_user;not null" json:"id_user"`
	User           User         `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user"`
	Deskripsi      string       `gorm:"column:deskripsi;type:text;not null" json:"deskripsi"`
	Status         string       `gorm:"column:status;type:enum('butuh tindakan','sedang diperbaiki','selesai');default:'butuh tindakan';not null" json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}