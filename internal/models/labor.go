package models

import "time"

type Labor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JurusanID uint      `gorm:"column:id_jurusan;not null" json:"id_jurusan"`
	Jurusan   Jurusan   `gorm:"foreignKey:JurusanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"jurusan"`
	Labor     string    `gorm:"size:100;not null" json:"labor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

