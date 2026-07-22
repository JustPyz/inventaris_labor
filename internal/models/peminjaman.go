package models

import "time"

type Peminjaman struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ItemInstanceID uint         `gorm:"column:id_item_instance;not null" json:"id_item_instance"`
	ItemInstance   ItemInstance `gorm:"foreignKey:ItemInstanceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"item_instance"`
	NamaPeminjam   string       `gorm:"column:nama_peminjam;size:150;not null" json:"nama_peminjam"`
	NomorTelepon   string       `gorm:"column:nomor_telepon;size:20;not null" json:"nomor_telepon"`
	TanggalPinjam  time.Time    `gorm:"column:tanggal_pinjam;not null" json:"tanggal_pinjam"`
	TanggalKembali time.Time    `gorm:"column:tanggal_kembali;not null" json:"tanggal_kembali"`
	Status         string       `gorm:"size:50;not null" json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
