package models

import "time"

type Perangkat struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NamaPerangkat string    `gorm:"column:nama_perangkat;size:150;not null" json:"nama_perangkat"`
	KategoriID    uint      `gorm:"column:kategori_id;not null" json:"kategori_id"`
	Kategori      Kategori  `gorm:"foreignKey:KategoriID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	JurusanID     uint      `gorm:"column:id_jurusan;not null" json:"id_jurusan"`
	Jurusan       Jurusan   `gorm:"foreignKey:JurusanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	LaborID       uint      `gorm:"column:id_labor;not null" json:"id_labor"`
	Labor         Labor     `gorm:"foreignKey:LaborID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Deskripsi     string    `gorm:"column:deskripsi;size:256;default:'tidak ada deskripsi'" json:"deskripsi"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
