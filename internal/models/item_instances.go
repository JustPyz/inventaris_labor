package models

import "time"

type ItemInstance struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PerangkatID uint      `gorm:"column:id_perangkat;not null" json:"id_perangkat"`
	Perangkat   Perangkat `gorm:"foreignKey:PerangkatID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"perangkat"`
	KodeAsset   string    `gorm:"column:kode_asset;size:100;not null;uniqueIndex" json:"kode_asset"`
	Status      string    `gorm:"size:50;not null" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}