package models

import "time"

type Perbaikan struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	KerusakanID        uint      `gorm:"column:id_kerusakan;not null" json:"id_kerusakan"`
	Kerusakan          Kerusakan `gorm:"foreignKey:KerusakanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"kerusakan"`
	UserID             uint      `gorm:"column:id_user;not null" json:"id_user"` // Teknisi / Kabeng
	User               User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user"`
	DeskripsiPerbaikan string    `gorm:"column:deskripsi_perbaikan;type:text;not null" json:"deskripsi_perbaikan"`
	Biaya              int64     `gorm:"column:biaya;not null;default:0" json:"biaya"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
