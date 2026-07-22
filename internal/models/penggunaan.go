package models

import "time"

type Penggunaan struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              uint      `gorm:"column:id_user;not null" json:"id_user"`
	User                User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	NamaPengguna        string    `gorm:"column:nama_pengguna;size:64;not null" json:"nama_pengguna"`
	LaborID             uint      `gorm:"column:id_labor;not null" json:"id_labor"`
	Labor               Labor     `gorm:"foreignKey:LaborID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"labor"`
	KelasID             uint      `gorm:"column:id_kelas;not null" json:"id_kelas"`
	Kelas               Kelas     `gorm:"foreignKey:KelasID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kelas"`
	JamPelajaranMulai   uint      `gorm:"column:jam_pelajaran_mulai;not null" json:"jam_pelajaran_mulai"`
	JamPelajaranSelesai uint      `gorm:"column:jam_pelajaran_selesai;not null" json:"jam_pelajaran_selesai"`
	CreatedAt           time.Time `json:"created_at"`
}
