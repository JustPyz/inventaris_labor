package models

import "time"

type Kelas struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JurusanID uint      `gorm:"column:id_jurusan;not null" json:"id_jurusan"`
	Jurusan   Jurusan   `gorm:"foreignKey:JurusanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"jurusan"`
	Kelas     string    `gorm:"size:100;not null" json:"kelas"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

