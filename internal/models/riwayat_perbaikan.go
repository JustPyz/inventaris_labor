package models

import "time"

// RiwayatPerbaikan adalah tabel audit independen yang menyimpan snapshot
// data lengkap pada saat setiap perbaikan selesai dicatat.
//
// Tabel ini TIDAK memiliki foreign key ke tabel manapun, sehingga:
//   - Data tetap ada meskipun perbaikan/kerusakan/item aslinya dihapus
//   - Bersifat append-only (hanya INSERT, tidak boleh UPDATE/DELETE via API)
//   - Semua field adalah nilai snapshot, bukan referensi
type RiwayatPerbaikan struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Referensi lunak (soft reference) — bukan FK, hanya untuk keperluan tracing
	PerbaikanID uint `gorm:"column:perbaikan_id;not null" json:"perbaikan_id"`
	KerusakanID uint `gorm:"column:kerusakan_id;not null" json:"kerusakan_id"`

	// Snapshot data item
	KodeAsset     string `gorm:"column:kode_asset;size:100;not null" json:"kode_asset"`
	NamaPerangkat string `gorm:"column:nama_perangkat;size:150;not null" json:"nama_perangkat"`

	// Snapshot data kerusakan
	DeskripsiKerusakan string `gorm:"column:deskripsi_kerusakan;type:text;not null" json:"deskripsi_kerusakan"`

	// Snapshot data perbaikan
	DeskripsiPerbaikan string `gorm:"column:deskripsi_perbaikan;type:text;not null" json:"deskripsi_perbaikan"`
	Biaya              int64  `gorm:"column:biaya;not null;default:0" json:"biaya"`

	// Snapshot data teknisi
	NamaTeknisi string `gorm:"column:nama_teknisi;size:64;not null" json:"nama_teknisi"`

	// Waktu perbaikan dicatat (diisi manual, bukan auto by GORM)
	TanggalPerbaikan time.Time `gorm:"column:tanggal_perbaikan;not null" json:"tanggal_perbaikan"`
}
