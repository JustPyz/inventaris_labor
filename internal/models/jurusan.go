package models

import "time"

type Jurusan struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	NamaJurusan string    `gorm:"size:150;not null" json:"nama_jurusan"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
